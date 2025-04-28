package routes

import (
	"github.com/acavaka/shortlinker/internal/handlers"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetURLHandler(w, r)
	case http.MethodPost:
		SaveURLHandler(w, r)
	default:
		handlers.RespondMethodNotAllowed(w, []string{"POST", "GET"})
	}
}
