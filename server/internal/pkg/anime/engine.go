package anime

import (
	"AnimeSearch/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func (handler *DBMaintainer) saveEdgeToCache(key string, value string, forward bool) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
	defer cancel()
	expTime := 600 * time.Second
	if forward {
		err = handler.rdb.Set(ctx, key, value, expTime).Err()
		if err != nil {
			return err
		}
		err = handler.rdb.Set(ctx, "-"+value, key, expTime).Err()
	} else {
		err = handler.rdb.Set(ctx, "-"+key, value, expTime).Err()
		if err != nil {
			return err
		}
		err = handler.rdb.Set(ctx, value, key, expTime).Err()
	}
	return
}

func (handler *DBMaintainer) getEdgeFromCache(key string, forward bool) (newKey string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
	defer cancel()
	defer func() {
		if newKey == "" {
			newKey = generateRandomHash(32)
		}
	}()
	if !forward {
		key = "-" + key
	}
	newKey = string(handler.rdb.Get(ctx, key).Val())
	return
}

func (handler *DBMaintainer) GetPreviousAnime(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Previous anime call received")
	key := chi.URLParam(r, "current_hash")
	if key == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	// Uno reverso card
	/*
		⣰⣾⣿⣿⠿⠿⢿⣿⣿⣿⣿⢿⣿⡿⣿⣿⣿⣿⣿⣿⣷⣆
		⣿⣟⣷⠋⢁⣀⣿⣟⣾⡷⣿⣻⡷⠛⠛⣊⣁⣉⠚⠳⣿⣻
		⣿⣿⠻⠂⢀⣿⣿⢯⣿⡽⠛⣀⣴⣶⣿⢿⣟⣿⣟⣦⡈⢻
		⣿⣯⣤⣴⣿⣟⣯⠟⢁⣴⣾⢿⣻⣽⣾⢿⣯⣿⢾⣟⣷⡀
		⣿⣿⣻⡿⣽⡟⢁⣴⡿⣯⡿⣟⣿⣽⣾⢿⣻⣾⢿⣯⡿⡇
		⣿⣯⣷⣿⠋⣠⣾⢿⣽⡿⣽⣿⣟⣾⣯⣿⢿⣽⣿⣳⣿⡇
		⣿⣯⡷⠁⣼⡿⣽⡿⣯⣿⣿⡦⠀⠀⠀⣸⣿⣻⣾⣟⣷⡇
		⣿⡿⢀⣼⣿⣻⣿⣽⣿⠟⠁⠀⠀⣠⠀⣿⡿⣽⡷⣿⡻⢀
		⣿⠀⣼⣿⣳⣿⣳⣿⣏⠀⢀⡴⠚⣿⣷⣿⢿⡿⣽⣟⠇⣸
		⡏⢰⣿⣯⣿⣻⣿⢿⣿⡤⠚⠁⠀⣿⣿⣟⣾⡷⣿⡟⢠⣿
		⢀⣾⣿⢾⣻⣽⣿⠀⠋⠀⠀⢀⣴⣿⣯⣿⢯⣿⡗⢁⣾⣿
		⢸⣿⢾⣿⣻⣿⡇⠀⠀⠀⠺⣿⣿⣻⣾⣟⣿⡝⢠⣿⣟⣿
		⢸⣿⣻⣽⡿⣿⣿⣿⣿⣿⣿⡿⣷⣿⣳⣿⠋⣴⣿⣻⣾⢿
		⢸⣿⣽⡷⣿⣻⣽⣯⣿⣽⡷⣿⣟⣾⠟⣡⡾⣟⣾⡿⣽⣿
		⠘⣿⢾⣟⣿⣽⡷⣿⡾⣯⣿⣻⠝⣡⣾⢿⣽⣿⡟⠛⢻⣿
		⣧⠈⠻⣿⣽⣾⣟⣿⡽⠿⢃⣥⣾⣟⣯⣿⢿⠊⠀⣤⣼⣿
		⣿⣿⣦⣤⣉⣁⣭⣥⣴⣾⣿⣻⣽⣾⢿⣽⠛⠃⢀⣼⣿⢿
		⠻⣷⣟⣯⣿⢯⣿⣽⢯⣷⣟⣯⣿⢾⡿⣿⣤⣤⣿⣿⢯⠟
		⠀⠀⠈⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
	*/

	previousKey := handler.getEdgeFromCache(key, false)
	err := handler.saveEdgeToCache(key, previousKey, false)
	if err != nil {
		fmt.Println("Unexpected error occurred: ", err)
	}
	handler.GetAnimeByKey(w, r, previousKey)
}

func (handler *DBMaintainer) GetNextAnime(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Next anime call received")
	key := chi.URLParam(r, "current_hash")
	if key == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	nextKey := handler.getEdgeFromCache(key, true)
	err := handler.saveEdgeToCache(key, nextKey, true)
	if err != nil {
		fmt.Println("Unexpected error occurred: ", err)
	}
	handler.GetAnimeByKey(w, r, nextKey)
}

func (handler *DBMaintainer) GetAnimeByKey(w http.ResponseWriter, r *http.Request, key string) {
	fmt.Println("Get anime call received")
	w.Header().Set("Content-Type", "application/json")
	animeID := hashToID(key, handler.getAnimeCount())
	var anime struct {
		CurrentHash string `json:"current_hash"`
		UrlPath     string `json:"url_path"`
		models.Anime
	}
	anime.CurrentHash = key
	var list []models.Anime
	handler.db.Offset(animeID).Limit(1).Find(&list)
	anime.Anime = list[0]
	anime.UrlPath = filepath.Join("/media/videos/", anime.AnimeHash, "cut.mp4")
	if anime.Id == 0 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	anime.Id = 0
	err := json.NewEncoder(w).Encode(anime)
	if err != nil {
		fmt.Println("Error encoding interview:", err)
	}
}

func (handler *DBMaintainer) GetAnime(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get anime call received")
	w.Header().Set("Content-Type", "application/json")
	key := chi.URLParam(r, "current_hash")
	if key == "" || len(key) < 8 {
		key = generateRandomHash(32)
	}
	animeID := hashToID(key, handler.getAnimeCount())
	var anime struct {
		CurrentHash string `json:"current_hash"`
		UrlPath     string `json:"url_path"`
		models.Anime
	}
	anime.CurrentHash = key
	var list []models.Anime
	handler.db.Offset(animeID).Limit(1).Find(&list)
	anime.Anime = list[0]
	anime.UrlPath = filepath.Join("/media/videos/", anime.AnimeHash, "cut.mp4")
	if anime.Id == 0 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	anime.Id = 0
	err := json.NewEncoder(w).Encode(anime)
	if err != nil {
		fmt.Println("Error encoding interview:", err)
	}
}

func createDirectory(path string) error {
	fmt.Println("Making new dir: ", path)
	return os.MkdirAll(path, os.ModePerm)
}

func (handler *DBMaintainer) CreateAnime(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var anime models.Anime
	err := decoder.Decode(&anime)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if anime.Name == "" || anime.OriginName == "" {
		http.Error(w, "Fields name and origin_name should not be empty", http.StatusBadRequest)
		return
	}
	anime.AnimeHash = generateRandomHash(16)
	handler.db.Create(&anime)
	dirPath := filepath.Join("./media/videos", anime.AnimeHash)

	err = createDirectory(dirPath)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}
	http.Error(w, "Anime created with id: "+anime.AnimeHash, http.StatusCreated)
}

func (handler *DBMaintainer) UploadAnime(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1 << 27) // Limit your input!
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form data: %v", err), http.StatusInternalServerError)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving the file: %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	hash := r.FormValue("anime_hash")

	destPath := filepath.Join("./media/videos", hash, "cut.mp4")
	destFile, err := os.Create(destPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating destination file: %v", err), http.StatusInternalServerError)
		return
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving file: %v", err), http.StatusInternalServerError)
		return
	}
	http.Error(w, http.StatusText(http.StatusCreated), http.StatusCreated)
}

func (handler *DBMaintainer) DeleteAnime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var anime models.Anime
	err := decoder.Decode(&anime)
	handler.db.First(&anime)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	oldDirPath := filepath.Join("./media/videos", anime.AnimeHash)
	delDirPath := filepath.Join("./media/videos", "deleted_"+anime.Name+"_"+generateRandomHash(10))
	err = os.Rename(oldDirPath, delDirPath)
	if err != nil {
		fmt.Println("Unexpected error while hiding old file: ", err)
		return
	}
	handler.db.Delete(&anime)
	err = json.NewEncoder(w).Encode(&anime)
	if err != nil {
		fmt.Println("Unexpected error while hiding anime: ", err)
	}
}