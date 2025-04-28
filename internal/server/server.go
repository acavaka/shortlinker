package server

import (
	"github.com/acavaka/shortlinker/internal/handlers/routes"
	"log"
	"net/http"
)

func RunServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", routes.RootHandler)

	log.Println("Server started on http://localhost:8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
