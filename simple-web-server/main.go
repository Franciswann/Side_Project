package main

import (
	"log"
	"net/http"
)

// main starts the HTTP server on port 8080
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", HelloHandler)
	mux.HandleFunc("/greet/", GreetHandler)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
