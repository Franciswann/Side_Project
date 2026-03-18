package main

import (
	"encoding/json"
	"net/http"
)

type Music struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
}

var musics = make(map[int]Music)

var musicList []Music

func init() {
	musics[1] = Music{
		Id:     1,
		Title:  "song_1",
		Artist: "Ed Sheeran",
	}
	musics[2] = Music{
		Id:     2,
		Title:  "song_2",
		Artist: "Daniel Caesar",
	}
	musics[3] = Music{
		Id:     3,
		Title:  "song_3",
		Artist: "De",
	}
}

func ListMusics(w http.ResponseWriter, r *http.Request) {

	for _, m := range musics {
		musicList = append(musicList, m)
	}

	data, err := json.Marshal(musicList)
	if err != nil {

	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
