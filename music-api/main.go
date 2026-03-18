package main

import (
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/musics", ListMusics)

	log.Println("Running...")

	http.ListenAndServe(":8080", nil)
}
