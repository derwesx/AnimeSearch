package dbha

import (
	"AnimeSearch/internal/pkg/log"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func ConnectRedis(log *log.Logger, cfg *viper.Viper) *redis.Client {
	log.Info("Connecting to Redis")
	defer func() { log.Info("Redis connected") }()
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:" + cfg.GetString("database.redis.port"),
		Password: cfg.GetString("database.redis.pass"),
		DB:       cfg.GetInt("database.redis.db"),
	})
}
