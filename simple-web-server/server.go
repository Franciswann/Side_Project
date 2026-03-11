package main

import (
	"net/http"
)

// HelloHandler responds with "Hello World" and status 200
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
