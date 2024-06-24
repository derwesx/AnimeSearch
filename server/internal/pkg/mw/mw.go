package mw

import (
	"AnimeSearch/internal/pkg/log"
	"github.com/spf13/viper"
)

type Middleware struct {
	config *viper.Viper
}

func New(log *log.Logger, config *viper.Viper) *Middleware {
	log.Info("Creating Middleware")
	defer func() { log.Info("Middleware created") }()
	return &Middleware{config: config}
}
