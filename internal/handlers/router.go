package handlers

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	mw "github.com/acavaka/shortlinker/internal/middleware"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

type URLSaverGetter interface {
	URLSaver
	URLGetter
}

func NewRouter(sg URLSaverGetter, baseURL string, logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()

	// All middleware must be defined before routes
	router.Use(chimiddleware.Recoverer)
	router.Use(mw.GzipMiddleware(logger))
	router.Use(mw.LoggerMiddleware(logger))

	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", GetHandler(sg, baseURL, logger))
		r.Post("/", SaveHandler(sg, baseURL, logger))
	})

	router.Post("/api/shorten", ShortenHandler(sg, baseURL, logger))

	return router
}
