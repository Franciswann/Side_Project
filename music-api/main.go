package main

import (
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("GET /musics", ListMusics)
	mux.HandleFunc("POST /musics", CreateMusic)

	// http.HandleFunc("/musics", ListMusics)

	log.Println("Running...")

	// ListenAndServe uses the configurable mux
	http.ListenAndServe(":8080", mux)
}
