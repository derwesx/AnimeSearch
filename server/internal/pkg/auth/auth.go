package auth

import (
	"AnimeSearch/internal/pkg/log"
	"AnimeSearch/internal/pkg/random"
	"AnimeSearch/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
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

	handler.db.Where("id = ?", (*claims)["Issuer"]).Select("id", "name", "email").First(&user)
	handler.logger.Debug(fmt.Sprint("User authenticated:", user))

	return user
}

func (handler *Auth) Authenticate(w http.ResponseWriter, r *http.Request) {
	var user = handler.IsAuthenticated(r)
	if user == nil {
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
	var user models.User
	err := decoder.Decode(&user)
	if err != nil {
		handler.logger.Debug("Login failed. Could not decode body.")
		http.Error(w, "Could not decode body.", http.StatusBadRequest)
		return
	}
	if user.Email == "" || user.Password == "" {
		handler.logger.Debug("Login failed. User not found.")
		http.Error(w, "User not found.", http.StatusBadRequest)
		return
	}
	var oldUser models.User
	handler.db.Where("email = ?", user.Email).First(&oldUser)
	if oldUser.Id == 0 {
		handler.logger.Debug("Login failed. User not found.")
		http.Error(w, "User not found.", http.StatusBadRequest)
		return
	}
	handler.logger.Debug("User logging: ", zap.String("params:", user.Password+" "+string(oldUser.Salt)+" "+oldUser.Password))
	if hashPassword(user.Password, oldUser.Salt) == oldUser.Password {
		handler.logger.Debug("Login succeeded.")
		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"Issuer":    strconv.Itoa(oldUser.Id),
			"ExpiresAt": time.Now().Add(time.Hour * 72).Unix(),
			"IsAdmin":   oldUser.IsAdmin,
		})
		token, err := claims.SignedString([]byte(handler.config.GetString("security.jwt.secret")))
		if err != nil {
			handler.logger.Info(fmt.Sprint("Unexpected error occurred:", err))
			return
		}
		cookie := http.Cookie{
			Name:     "jwt",
			Value:    token,
			Expires:  time.Now().Add(time.Hour * 72),
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)

		oldUser.Password = ""
		oldUser.Salt = make([]byte, 0)
		err = json.NewEncoder(w).Encode(oldUser)
		if err != nil {
			handler.logger.Info(fmt.Sprint("Login | Unexpected error occurred:", err))
		}
	} else {
		handler.logger.Debug("Login failed. Email or password incorrect.")
		http.Error(w, "Login failed.", http.StatusBadRequest)
		return
	}
}

func (handler *Auth) Logout(w http.ResponseWriter, r *http.Request) {
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
	var user models.User
	err := decoder.Decode(&user)
	if err != nil {
		handler.logger.Debug("Registration failed. Could not decode body.")
		http.Error(w, "Registration failed.", http.StatusBadRequest)
		return
	}
	user.Salt = random.GenerateSalt(handler.config.GetInt("security.auth.salt_size"))
	user.Password = hashPassword(user.Password, user.Salt)
	user.IsAdmin = false
	user.Favourites = make(pq.Int64Array, 0)
	handler.logger.Debug(fmt.Sprint("Registering user:", user))
	handler.db.Create(&user)

	// Pls fix it ASAP
	user.Salt = make([]byte, 0)
	user.Password = ""
	user.IsAdmin = false

	if user.Id == 0 {
		handler.logger.Debug("Registration failed.")
		http.Error(w, "Registration failed.", http.StatusBadRequest)
		return
	}
	handler.logger.Debug("User registered successfully.")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		handler.logger.Fatal(fmt.Sprint("Register | Unexpected error occurred:", err))
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
