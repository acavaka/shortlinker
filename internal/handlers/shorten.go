package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/acavaka/shortlinker/internal/models"
	"go.uber.org/zap"
)

type URLSaver interface {
	SaveURL(url string) string
}

func ShortenHandler(saver URLSaver, baseURL string, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.Request
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			logger.Error("failed to decode request body", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		short := saver.SaveURL(req.URL)
		resultURL, err := url.JoinPath(baseURL, short)
		if err != nil {
			logger.Error("failed to join path to get result URL", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		resp := models.Response{
			Result: resultURL,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		enc := json.NewEncoder(w)
		if err = enc.Encode(resp); err != nil {
			logger.Error("failed to encode response", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
