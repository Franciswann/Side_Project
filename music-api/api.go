package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
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

// ListMusics list all the musics
func ListMusics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	val, err := rdb.Get(ctx, "all_musics").Result()

	// if key does not exist, extract data from database and store data into Redis
	if err == redis.Nil {

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

		// store cache data in redis with 5-minute expiration
		err = rdb.Set(ctx, "all_musics", data, 5*time.Minute).Err()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(data); err != nil {
			log.Printf("Failed to write response: %v", err)
			return
		}

	} else {
		// Return cached data
		w.WriteHeader(http.StatusOK)

		if _, err = w.Write([]byte(val)); err != nil {
			log.Printf("Failed to write response: %v", err)
			return
		}
	}
}

// CreateMusic create new music
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

		// Music already exists
		if err == nil {
			w.WriteHeader(http.StatusConflict)
			if _, err := w.Write([]byte(fmt.Sprintf("Music with %s and %s already exists", newMusic.Title, newMusic.Artist))); err != nil {
				log.Printf("Failed to write response: %v", err)
			}
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

		// Invalidate cache - data has changed
		err = rdb.Del(ctx, "all_musics").Err()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed to invalidate cache: %v", err)))
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(data)
	}
}

// GetMusic fetch specific music
func GetMusic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var searchedMusic Music

	// Extract id from URL path: /musics/{id}
	path := path.Base(r.URL.Path)
	id, _ := strconv.Atoi(path)

	// Get specific cached data from Redis
	val, err := rdb.Get(ctx, path).Result()

	// Music doesn't exist in Redis, extract data from Postgres and cache into Redis
	if err == redis.Nil {
		err := db.QueryRow(`SELECT id, title, artist FROM musics WHERE id=$1;`, id).Scan(&searchedMusic.Id, &searchedMusic.Title, &searchedMusic.Artist)
		switch {
		// couldn't find the music
		case err == sql.ErrNoRows:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("no music with id %d", id)))
			return
		case err != nil:
			log.Printf("query error: %v\n", err)
			return
		// successfully found the music
		default:
			data, err := json.Marshal(searchedMusic)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// store cache data in redis with 3-minute expiration
			err = rdb.Set(ctx, path, data, 3*time.Minute).Err()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(data)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	// Return cached data
	if _, err = w.Write([]byte(val)); err != nil {
		log.Printf("Failed to write response: %v", err)
		return
	}
}

// DeleteMusic deletes a specific music
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

	// rows: the number of rows affected by DELETE
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

	// Invalidate cache - data has changed
	err = rdb.Del(ctx, "all_musics", path).Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Failed to invalidate cache: %v", err)))
	}

	// Successfully deleted
	w.WriteHeader(http.StatusNoContent)
}

// UpdateMusic update specific music by ID
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

		// Invalidate cache - data has changed
		err = rdb.Del(ctx, "all_musics").Err()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed to invalidate cache: %v", err)))
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}
