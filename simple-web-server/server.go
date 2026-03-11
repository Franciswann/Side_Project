package main

import (
	"fmt"
	"net/http"
	"strings"
)

// HelloHandler responds with "Hello World" and status 200
// type conversion []byte()
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

// GreetHandler reads the name from URL path and responds "Hi, {name}!"
// type conversion []byte()
func GreetHandler(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/greet/")
	fmt.Fprintf(w, "Hi, %s!", name)
}
