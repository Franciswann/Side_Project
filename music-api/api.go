package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
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
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func CreateMusic(w http.ResponseWriter, r *http.Request) {

	var newMusic Music

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&newMusic); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Marshaling failed."))
		return

	} else {
		maxID := 0
		for i, _ := range musics {
			if i > maxID {
				maxID = i
			}
		}
		newID := maxID + 1
		newMusic.Id = newID
		musics[newID] = newMusic

		data, err := json.Marshal(newMusic)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(data)
	}

}

func GetMusic(w http.ResponseWriter, r *http.Request) {

	path := path.Base(r.URL.Path)
	id, _ := strconv.Atoi(path)

	// music, ok:= musics[id] --- Comma-ok idiom

	// if exist(ok = true) then return music info and status 200
	// if not exist return "Music not found" and status 400
	if music, ok := musics[id]; ok {
		music, err := json.Marshal(music)
		if err != nil {
			fmt.Println("Marshaling failed")
		}

		w.WriteHeader(http.StatusOK)
		w.Write(music)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Music not found"))
	}
}
