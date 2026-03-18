package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Music struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
}

var musics = make(map[int]Music)

func init() {
	musics[1] = Music{
		Id:     1,
		Title:  "Perfect",
		Artist: "Ed Sheeran",
	}
	musics[2] = Music{
		Id:     2,
		Title:  "Always",
		Artist: "Daniel Caesar",
	}
	musics[3] = Music{
		Id:     3,
		Title:  "Die For You",
		Artist: "Joji",
	}
}

func ListMusics(w http.ResponseWriter, r *http.Request) {

	musicList := []Music{}

	for _, m := range musics {
		musicList = append(musicList, m)
	}

	data, err := json.Marshal(musicList)
	// neet to fix, make it show on curl HTTP, and show statusCode
	if err != nil {
		fmt.Fprintf(w, string(http.StatusInternalServerError))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func CreateMusic(w http.ResponseWriter, r *http.Request) {
	// JSON body

	// var newMusic Music
	// const newSong = `{"title":"Blinding Lights","artist":"The Weekend"}`

	// dec := json.NewDecoder(strings.NewReader(newSong)).Decode(&newMusic)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("POST received"))

}
