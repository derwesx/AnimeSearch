package main

import (
	"AnimeSearch/internal/app/AnimeSearch"
	"AnimeSearch/internal/pkg/config"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	cfg := config.NewConfig()
	r := chi.NewRouter()
	AnimeSearch.SetupApp(cfg, r)
	log.Fatal(http.ListenAndServe(":"+cfg.GetString("http.port"), r))
}
