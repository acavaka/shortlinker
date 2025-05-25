package handlers

import (
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"

	"github.com/acavaka/shortlinker/internal/logger"
	"github.com/acavaka/shortlinker/internal/service"
)

func GetHandler(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		short := chi.URLParam(r, "id")
		long, err := svc.GetURL(short)
		if err != nil {
			logger.Error("failed to get URL", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		origin, err := url.JoinPath(svc.BaseURL, short)
		if err != nil {
			logger.Error("failed to join path to get redirect URL", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", origin)
		http.Redirect(w, r, long, http.StatusTemporaryRedirect)
	}
}
