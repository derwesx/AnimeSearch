package dbha

import (
	"AnimeSearch/internal/pkg/log"
	"AnimeSearch/models"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgres(log *log.Logger, cfg *viper.Viper) *gorm.DB {
	log.Info("Connecting to Postgres")
	defer func() { log.Info("Postgres connected") }()
	user := cfg.GetString("database.postgres.user")
	pass := cfg.GetString("database.postgres.pass")
	port := cfg.GetString("database.postgres.port")
	database := cfg.GetString("database.postgres.db")
	dbURL := "postgres://" + user + ":" + pass + "@localhost:" + port + "/" + database + "?sslmode=disable"
	log.Info("Connecting to PostgreSQL", zap.String("url", dbURL))
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dbURL,
	}), &gorm.Config{})

	if err != nil {
		log.Fatal(fmt.Sprint("Couldn't connect Postgres", err))
	}

	log.Info("Migrating Postgres")
	err = db.AutoMigrate(&models.User{}, &models.Anime{})
	if err != nil {
		log.Fatal(fmt.Sprint("Couldn't migrate Postgres", err))
	}
	log.Info("Postgres successfully migrated")

	return db
}
