package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// HelloHandler responds with JSON {"message": "Hello World"} and correct Content-Type
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "Hello World"}
	jsonBytes, _ := json.Marshal(response) //ignore error for this simple case
	w.Write(jsonBytes)
}

// GreetHandler reads the name from URL path and responds "Hi, {name}!"
func GreetHandler(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/greet/")
	fmt.Fprintf(w, "Hi, %s!", name)
}
