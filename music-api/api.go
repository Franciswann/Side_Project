package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
			w.Write([]byte(fmt.Sprintf("Music with %s and %s already exists", newMusic.Title, newMusic.Artist)))
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

	var searchedMusic Music

	// Extract id from URL path: /musics/{id}
	path := path.Base(r.URL.Path)
	id, _ := strconv.Atoi(path)

	err := db.QueryRow(`SELECT id, title, artist FROM musics WHERE id=$1;`, id).Scan(&searchedMusic.Id, &searchedMusic.Title, &searchedMusic.Artist)
	switch {
	// couldn't find the music
	case err == sql.ErrNoRows:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("no music with id %d", id)))
	case err != nil:
		log.Printf("query error: %v\n", err)
	// successfully found the music
	default:
		data, err := json.Marshal(searchedMusic)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

// Delete specific music
func DeleteMusic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract id from URL path: /musics/{id}
	path := path.Base(r.URL.Path)
	id, _ := strconv.Atoi(path)

	result, err := db.Exec(`DELETE FROM musics WHERE id=$1;`, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Delete error: %v", err)
		return
	}

	// rows: the number of rows afftected by DELETE
	rows, err := result.RowsAffected()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("RowsAffected error: %v", err)
		return
	}
	// music not found
	if rows == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("no music with id %d", id)))
		return
	}
	// Successfully deleted
	w.WriteHeader(http.StatusNoContent)
}

// Update specific music by ID
func UpdateMusic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var updateMusic Music

	// Parse and decode JSON request body
	if err := json.NewDecoder(r.Body).Decode(&updateMusic); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Marshaling failed."))
		return
	} else {
		// Extract music ID from URL path: /musics/{id}
		path := path.Base(r.URL.Path)
		id, _ := strconv.Atoi(path)

		// Update music record and return updated data
		updateQuery := `UPDATE musics
						SET title=$1, artist=$2
						WHERE id=$3
						RETURNING id, title, artist;`
		err := db.QueryRow(updateQuery, updateMusic.Title, updateMusic.Artist, id).Scan(&updateMusic.Id, &updateMusic.Title, &updateMusic.Artist)

		// Music not found
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("no music with id: %d", id)))
			return
		}

		// Database error
		if err != nil {
			log.Printf("Update error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Return updated music data
		data, err := json.Marshal(updateMusic)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}
