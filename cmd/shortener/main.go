package main

import (
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

var (
	urlStorage = make(map[string]string)
	storageMu  sync.RWMutex
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", routeHandler)

	log.Println("Server started on http://localhost:8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

}

func routeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleShortenURL(w, r)
	case http.MethodGet:
		handleRedirectToOriginalURL(w, r)
	default:
		respondMethodNotAllowed(w)
	}
}

func handleShortenURL(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	longURL := strings.TrimSpace(string(body))
	if longURL == "" {
		http.Error(w, "empty URL is not allowed", http.StatusBadRequest)
		return
	}

	shortURL := "EwHXdJfB"

	storageMu.Lock()
	urlStorage[shortURL] = longURL
	storageMu.Unlock()

	responseURL := buildShortURL(r, shortURL)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(responseURL))
}

func handleRedirectToOriginalURL(w http.ResponseWriter, r *http.Request) {
	shortURL := strings.TrimPrefix(r.URL.Path, "/")
	if shortURL == "" {
		http.Error(w, "invalid short URL", http.StatusBadRequest)
		return
	}

	storageMu.RLock()
	longURL, exists := urlStorage[shortURL]
	storageMu.RUnlock()

	if !exists {
		http.Error(w, "short URL not found", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
}

func respondMethodNotAllowed(w http.ResponseWriter) {
	http.Error(w, "method not allowed, use GET or POST", http.StatusMethodNotAllowed)
}

func buildShortURL(r *http.Request, shortID string) string {
	return "http://" + r.Host + "/" + shortID
}
