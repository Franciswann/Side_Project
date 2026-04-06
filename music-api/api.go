package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"path"
	"strconv"

	_ "github.com/lib/pq"
)

type Music struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
}

var musics = make(map[int]Music)
var db *sql.DB

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

func initDB() error {
	// insert initial data to table 'musics' with conflict handling
	for _, music := range musics {
		query := `INSERT INTO musics (title, artist) 
				  VALUES ($1, $2) 
				  ON CONFLICT (title, artist) DO NOTHING;`
		_, err := db.Exec(query, music.Title, music.Artist)
		if err != nil {
			return err
		}
	}
	return nil
}

// List all the musics
func ListMusics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// delete this part

	// for _, m := range musics {
	// 	musicList = append(musicList, m)
	// }
	rows, err := db.Query("SELECT * FROM musics")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	musicList := []Music{}

	// reverse musics data table and append into musicList
	for rows.Next() {
		// musicList := []Music{}
		var musicInRows Music
		if err := rows.Scan(&musicInRows.Id, &musicInRows.Title, &musicInRows.Artist); err != nil {
			log.Println(err)
		}
		musicList = append(musicList, musicInRows)
	}

	data, err := json.Marshal(musicList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

// create new music
func CreateMusic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newMusic Music

	if err := json.NewDecoder(r.Body).Decode(&newMusic); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Marshaling failed."))
		return
	} else {
		// Check if music with same title and artist already exists
		var existingId int
		checkQuery := `SELECT id FROM musics WHERE title=$1 AND artist=$2`
		err := db.QueryRow(checkQuery, newMusic.Title, newMusic.Artist).Scan(&existingId)

		if err == nil {
			// Music already exists
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("Music with this title and artist already exists"))
			return
		}

		// INSERT without Id, it would be serial automatic generate
		query := `INSERT INTO musics (title, artist)
				  VALUES ($1, $2)
				  RETURNING id;`

		if err := db.QueryRow(query, newMusic.Title, newMusic.Artist).Scan(&newMusic.Id); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to create music"))
			return
		}

		data, err := json.Marshal(newMusic)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(data)
	}

}

// fetch specific music
func GetMusic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := path.Base(r.URL.Path)
	id, _ := strconv.Atoi(path)

	// music, ok:= musics[id] --- Comma-ok idiom

	// if exist(ok = true) then return music info and status 200
	// if not exist return "Music not found" and status 400
	if music, ok := musics[id]; ok {
		data, err := json.Marshal(music)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Invalid JSON format"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Music not found"))
	}
}

// Delete specific music
func DeleteMusic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := path.Base(r.URL.Path)
	id, _ := strconv.Atoi(path)

	// delete(map, key)
	if _, ok := musics[id]; ok {
		delete(musics, id)
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Music not found"))
	}
}

// Parsing request body and update musics
func UpdateMusic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var updateMusic Music

	// parsing and decoding json and store into upupdateMusic
	if err := json.NewDecoder(r.Body).Decode(&updateMusic); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Marshaling failed."))
		return
	} else {
		path := path.Base(r.URL.Path)
		id, _ := strconv.Atoi(path)

		if _, ok := musics[id]; ok {
			w.WriteHeader(http.StatusOK)
			updateMusic.Id = id
			musics[id] = updateMusic

			data, err := json.Marshal(updateMusic)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Write(data)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}
