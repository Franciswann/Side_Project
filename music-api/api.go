package main

import (
	"net/http"
)

type Music struct {
	Id     int64  `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
}

func ListMusics(w http.ResponseWriter, r *http.Request) {

}
