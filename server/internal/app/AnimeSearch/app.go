package AnimeSearch

import (
	"AnimeSearch/internal/pkg/anime"
	"AnimeSearch/internal/pkg/auth"
	"AnimeSearch/internal/pkg/dbha"
	"AnimeSearch/internal/pkg/log"
	"AnimeSearch/internal/pkg/mw"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"net/http"
	"os"
	"path/filepath"
)

var configGlobal *viper.Viper
var loggerGlobal *log.Logger

var maintainer *anime.DBMaintainer
var middleware *mw.Middleware
var authHandler *auth.Auth

func setupPostgres() *gorm.DB {
	db := dbha.ConnectPostgres(loggerGlobal, configGlobal)
	return db
}

func setupRedis() *redis.Client {
	rdb := dbha.ConnectRedis(loggerGlobal, configGlobal)
	return rdb
}

func setupMaintainer(db *gorm.DB, rdb *redis.Client) {
	maintainer = anime.NewMaintainer(db, rdb)
}

func setupMiddleware() {
	middleware = mw.New(loggerGlobal, configGlobal)
}

func setupAuth(db *gorm.DB, rdb *redis.Client) {
	authHandler = auth.New(loggerGlobal, configGlobal, db, rdb)
}

func setupRoutes(r *chi.Mux) {
	r.Use(mw.CORSMiddleware)
	// <-- Auth -->
	// Input: Name, Email, Password
	r.Post("/api/register", authHandler.Register)
	// Input: Email, Password
	r.Post("/api/login", authHandler.Login)
	// Should be authenticated
	r.Get("/api/logout", authHandler.Logout)
	// Should be authenticated. Gets current user
	r.Get("/api/user", authHandler.Authenticate)

	// select first_name from person order by first_name desc limit 1 offset 8;
	r.Get("/api/anime/next/{current_hash}", maintainer.GetNextAnime)
	r.Get("/api/anime/prev/{current_hash}", maintainer.GetPreviousAnime)
	r.Get("/api/anime/{current_hash}", maintainer.GetAnime)

	r.With(middleware.Admin).Post("/api/anime/create", maintainer.CreateAnime)
	r.With(middleware.Admin).Post("/api/anime/upload", maintainer.UploadAnime)
	r.With(middleware.Admin).Post("/api/anime/delete", maintainer.DeleteAnime)

	// Just a long request for getting video files. TODO make a module with it
	r.Get("/media/videos/{hash}/{episode}/{filename}", func(w http.ResponseWriter, r *http.Request) {
		hash := chi.URLParam(r, "hash")
		episode := chi.URLParam(r, "episode")
		filename := chi.URLParam(r, "filename")

		// Construct the absolute file path
		absFilePath := filepath.Join("./media/videos", hash, episode, filename)

		// Open the file
		file, err := os.Open(absFilePath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error opening file: %v", err), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Get file info (to obtain file size)
		fileInfo, err := file.Stat()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error getting file info: %v", err), http.StatusInternalServerError)
			return
		}

		// Set HTTP headers
		w.Header().Set("Content-Type", "video/mp4") // Adjust content type based on your file type
		w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
		w.Header().Set("Accept-Ranges", "bytes") // Enable support for range requests

		// Serve the file content using http.ServeContent
		http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
	})
}

func SetupApp(cfg *viper.Viper, r *chi.Mux) {
	configGlobal = cfg
	loggerGlobal = log.New(configGlobal)
	db := setupPostgres()
	rdb := setupRedis()
	setupMaintainer(db, rdb)
	setupMiddleware()
	setupAuth(db, rdb)
	setupRoutes(r)
}
