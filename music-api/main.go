package main

import (
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("GET /musics", ListMusics)
	// GetMusic handles GET requests to /musics/{id}
	// Example: GET /musics/4
	mux.HandleFunc("GET /musics/{id}", GetMusic)

	mux.HandleFunc("POST /musics", CreateMusic)

	log.Println("Running...")

	// ListenAndServe uses the configurable mux
	// http://localhost:8080/
	http.ListenAndServe(":8080", mux)
}
