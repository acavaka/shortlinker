package handlers

import (
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	mw "github.com/acavaka/shortlinker/internal/middleware"
	"github.com/acavaka/shortlinker/internal/service"
)

func NewRouter(svc *service.Service) *chi.Mux {
	router := chi.NewRouter()

	// All middleware must be defined before routes
	router.Use(chimiddleware.Recoverer)
	router.Use(mw.GzipMiddleware)
	router.Use(mw.Logger)

	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", GetHandler(svc))
		r.Post("/", SaveHandler(svc))
	})

	router.Post("/api/shorten", ShortenHandler(svc))

	return router
}
