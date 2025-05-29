package handlers

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"go.uber.org/zap"
)

func SaveHandler(svc URLSaver, baseURL string, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		long, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("failed to read body", zap.Error(err))
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		urlStr := strings.TrimSpace(string(long))
		if urlStr == "" {
			http.Error(w, "URL cannot be empty", http.StatusBadRequest)
			return
		}

		parsedURL, err := url.ParseRequestURI(urlStr)
		if err != nil {
			http.Error(w, "Invalid URL format", http.StatusBadRequest)
			return
		}

		if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
			http.Error(w, "URL must use http or https protocol", http.StatusBadRequest)
			return
		}

		short := svc.SaveURL(urlStr)

		resultURL, err := url.JoinPath(baseURL, short)
		if err != nil {
			logger.Error("failed to join path to get result URL", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(resultURL))
		if err != nil {
			logger.Error("failed to write the full URL response to client", zap.Error(err))
		}
	}
}
