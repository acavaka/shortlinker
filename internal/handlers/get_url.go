package handlers

import (
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type URLGetter interface {
	GetURL(shortURL string) (string, error)
}

func GetHandler(svc URLGetter, baseURL string, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		short := chi.URLParam(r, "id")
		long, err := svc.GetURL(short)
		if err != nil {
			logger.Error("failed to get URL", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		origin, err := url.JoinPath(baseURL, short)
		if err != nil {
			logger.Error("failed to join path to get redirect URL", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", origin)
		http.Redirect(w, r, long, http.StatusTemporaryRedirect)
	}
}
