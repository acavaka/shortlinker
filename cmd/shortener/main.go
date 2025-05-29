package main

import (
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/acavaka/shortlinker/internal/config"
	"github.com/acavaka/shortlinker/internal/handlers"
	"github.com/acavaka/shortlinker/internal/logger"
	"github.com/acavaka/shortlinker/internal/service"
	"github.com/acavaka/shortlinker/internal/storage"
	"go.uber.org/zap"
)

func normalizeAddress(addr string) string {
	if strings.Contains(addr, "[") {
		return addr
	}

	if strings.HasPrefix(addr, ":") {
		return addr
	}

	if addr == "localhost:8080" || addr == "127.0.0.1:8080" {
		return "[::1]:8080"
	}

	return addr
}

func main() {
	log := logger.Initialize()
	defer log.Sync()

	cfg := config.LoadConfig()

	normalizedAddr := normalizeAddress(cfg.Service.ServerAddress)

	var (
		db  storage.URLStorage
		err error
	)

	if cfg.Service.FileStoragePath == "" {
		db = storage.NewMemoryStorage(cfg)
	} else {
		db, err = storage.NewFileStorage(cfg, log)
		if err != nil {
			log.Fatal("failed to load storage", zap.Error(err))
		}
	}

	svc := &service.Service{DB: db, BaseURL: cfg.Service.BaseURL, FileStoragePath: cfg.Service.FileStoragePath}
	r := handlers.NewRouter(svc, cfg.Service.BaseURL, log)

	listener, err := net.Listen("tcp", normalizedAddr)
	if err != nil {
		log.Fatal("failed to create listener", zap.Error(err))
	}

	srv := &http.Server{
		Addr:    cfg.Service.ServerAddress,
		Handler: r,
	}

	log.Info("server started", zap.String("address", listener.Addr().String()))
	if err := srv.Serve(listener); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Error("unexpected error", zap.Error(err))
		}
	}
}
