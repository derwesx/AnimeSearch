package auth

import (
	"AnimeSearch/internal/pkg/log"
	"AnimeSearch/internal/pkg/random"
	"AnimeSearch/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type Auth struct {
	config *viper.Viper
	logger *log.Logger
	db     *gorm.DB
	rdb    *redis.Client
}

func (handler *Auth) IsAuthenticated(r *http.Request) (user *models.User) {
	cookie, err := r.Cookie("jwt")
	if errors.Is(err, http.ErrNoCookie) {
		return nil
	}

	token, err := jwt.ParseWithClaims(cookie.Value, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(handler.config.GetString("security.jwt.secret")), nil
	})

	if err != nil {
		return nil
	}

	claims := token.Claims.(*jwt.MapClaims)

	handler.db.Where("id = ?", (*claims)["Issuer"]).First(&user)
	handler.logger.Debug("User authenticated:")

	return user
}

func (handler *Auth) Authenticate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user = handler.IsAuthenticated(r)
	if user == nil || user.Id == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	err := json.NewEncoder(w).Encode(user)
	if err != nil {
		handler.logger.Fatal(fmt.Sprint("Unexpected error happened:", err))
	}
}

func (handler *Auth) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var user struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := decoder.Decode(&user)

	if err != nil {
		handler.logger.Debug("Login failed. Could not decode body.")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if user.Email == "" || user.Password == "" {
		handler.logger.Debug("Login failed. User not found.")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	var oldUser models.User
	handler.db.Where("email = ?", user.Email).First(&oldUser)
	if oldUser.Id == 0 {
		handler.logger.Debug("Login failed. User not found.")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if matchPassword(oldUser.Password, user.Password, oldUser.Salt) {
		handler.logger.Debug("Login succeeded.")

		// Creating new JWT token
		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"Issuer":    strconv.Itoa(oldUser.Id),
			"ExpiresAt": time.Now().Add(time.Hour * 72).Unix(),
			"IsAdmin":   oldUser.IsAdmin,
		})
		token, err := claims.SignedString([]byte(handler.config.GetString("security.jwt.secret")))
		if err != nil {
			handler.logger.Fatal(fmt.Sprint("Unexpected error occurred:", err))
			return
		}
		cookie := http.Cookie{
			Name:     "jwt",
			Value:    token,
			Expires:  time.Now().Add(time.Hour * 72),
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)

		err = json.NewEncoder(w).Encode(oldUser)
		if err != nil {
			handler.logger.Info(fmt.Sprint("Unexpected error occurred:", err))
		}
		return
	}

	http.Error(w, "Login failed. Email or password is incorrect.", http.StatusNotFound)
}

func (handler *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}

func (handler *Auth) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var user struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := decoder.Decode(&user)
	if err != nil {
		handler.logger.Debug("Registration failed. Could not decode body.")
		http.Error(w, "Registration failed.", http.StatusBadRequest)
		return
	}
	newUser := models.User{
		Name:  user.Name,
		Email: user.Email,
	}
	newUser.Salt = random.GenerateSalt(handler.config.GetInt("security.auth.salt_size"))
	newUser.Password = hashPassword(user.Password, newUser.Salt)
	handler.db.Create(&newUser)

	if newUser.Id == 0 {
		handler.logger.Debug("Registration failed. User not created.")
		http.Error(w, "Registration failed.", http.StatusBadRequest)
		return
	}
	handler.logger.Debug("User registered successfully.")
	err = json.NewEncoder(w).Encode(newUser)
	if err != nil {
		handler.logger.Fatal(fmt.Sprint("Unexpected error occurred:", err))
	}
}

func New(log *log.Logger, cfg *viper.Viper, db *gorm.DB, rdb *redis.Client) *Auth {
	log.Info("Initializing auth handler")
	defer func() { log.Info("Auth handler initialized") }()
	return &Auth{
		config: cfg,
		logger: log,
		db:     db,
		rdb:    rdb,
	}
}
