package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/acavaka/shortlinker/internal/config"
	"github.com/acavaka/shortlinker/internal/handlers"
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
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	cfg := config.LoadConfig()

	normalizedAddr := normalizeAddress(cfg.Server.ServerAddress)

	db := storage.LoadStorage()
	svc := &service.Service{
		DB:      db,
		BaseURL: cfg.Server.BaseURL,
	}
	r := handlers.NewRouter(svc, logger)

	listener, err := net.Listen("tcp", normalizedAddr)
	if err != nil {
		logger.Fatal("failed to create listener", zap.Error(err))
	}

	srv := &http.Server{
		Handler: r,
	}

	logger.Info("server started", zap.String("address", listener.Addr().String()))
	if err := srv.Serve(listener); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("unexpected error", zap.Error(err))
		}
	}
}
