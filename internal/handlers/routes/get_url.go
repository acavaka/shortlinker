package routes

import (
	"github.com/acavaka/shortlinker/internal/handlers"
	"github.com/acavaka/shortlinker/internal/storage"
	"net/http"
	"strings"
)

func GetURLHandler(w http.ResponseWriter, r *http.Request) {
	shortLink := strings.TrimPrefix(r.URL.Path, "/")
	handlers.SetHeadersHandler(w)

	mapper := storage.Mapper
	longLink, ok := mapper.Get(shortLink)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
	}

	originalURL := handlers.BuildShortURL(r, shortLink)
	handlers.RedirectToURL(w, r, longLink, originalURL)

	_, err := w.Write([]byte(longLink))
	if err != nil {
		return
	}
}
