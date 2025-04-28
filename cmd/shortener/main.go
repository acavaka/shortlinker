package main

import (
	"log"
	"net/http"

	"github.com/acavaka/shortlinker/internal/config"
	"github.com/acavaka/shortlinker/internal/handlers"
	"github.com/acavaka/shortlinker/internal/service"
	"github.com/acavaka/shortlinker/internal/storage"
)

func main() {
	cfg := config.LoadConfig()
	db := storage.LoadStorage()
	svc := &service.Service{DB: db, BaseURL: cfg.Server.BaseURL}
	r := handlers.NewRouter(svc)

	srv := &http.Server{
		Addr:    cfg.Server.ServerAddress,
		Handler: r,
	}

	log.Printf("server started on: %s", cfg.Server.ServerAddress)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("got unexpected error, details: %s", err)
	}
}
