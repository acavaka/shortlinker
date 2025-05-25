package handlers

import (
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/acavaka/shortlinker/internal/middleware"
	"github.com/acavaka/shortlinker/internal/service"
	"go.uber.org/zap"
)

func NewRouter(svc *service.Service, logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()

	// All middleware must be defined before routes
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.SetHeader("Content-Type", "text/plain; charset=utf-8"))
	router.Use(middleware.WithLogging(logger))

	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", GetHandler(svc))
		r.Post("/", SaveHandler(svc))
	})

	router.Route("/api", func(r chi.Router) {
		r.Post("/shorten", ShortenHandler(svc))
	})

	return router
}
