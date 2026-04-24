package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"music-api/internal/database"
	"music-api/internal/handler"
	"music-api/internal/model"

	"github.com/redis/go-redis/v9"
)

// Truncate data table and insert testing data
func setupTestDB(t *testing.T) {

	// Truncate table
	_, err := database.DB.Exec(`TRUNCATE TABLE musics RESTART IDENTITY;`)
	if err != nil {
		t.Errorf("Failed to truncate table: %v", err)
	}

	// Insert test data
	testMusic := []model.Music{
		{Title: "Always", Artist: "Daniel Caesar"},
		{Title: "Die For You", Artist: "Joji"},
	}

	for _, music := range testMusic {
		query := `INSERT INTO musics (title, artist)
				  VALUES ($1, $2)
				  ON CONFLICT (title, artist) DO NOTHING;`
		_, err := database.DB.Exec(query, music.Title, music.Artist)
		if err != nil {
			t.Errorf("Failed to insert music: %v", err)
		}
	}
}

func TestMain(m *testing.M) {
	// Initialize db connection
	connStr := "host=localhost port=5432 user=wanchaochun password=password dbname=music_db sslmode=disable"
	var err error
	database.DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.DB.Close()

	// Initialize redis connection
	if database.RDB == nil {
		database.RDB = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})

		_, err := database.RDB.Ping(database.Ctx).Result()
		if err != nil {
			log.Printf("connection failed: %s", err)
		} else {
			log.Printf("Connected to Redis successfully")
		}
		defer database.RDB.Close()
	}

	// Execute all tests
	code := m.Run()

	os.Exit(code)
}

func TestListMusics(t *testing.T) {
	// Setup test data
	setupTestDB(t)

	// Clear Redis cache
	database.RDB.Del(database.Ctx, "all_musics")

	req := httptest.NewRequest("GET", "/musics", nil)
	recorder := httptest.NewRecorder()

	// because HandlerFunc already implement ServeHTTP, so
	// ListMusics can seen as Handler
	handler.ListMusics(recorder, req)

	// Parsing response data
	var music []model.Music
	err := json.Unmarshal(recorder.Body.Bytes(), &music)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	// Test status code
	t.Run("Validate http status code", func(t *testing.T) {
		if status := recorder.Code; status != http.StatusOK {
			t.Errorf("got %v, want %v", status, http.StatusOK)
		}
	})

	// Test response data
	t.Run("Validate data amount", func(t *testing.T) {
		// check if the length == 2
		got := len(music)
		if got != 2 {
			t.Errorf("Expected 2 items, got %d", got)
		}
	})

	// Test specified data
	t.Run("Validate specific titles in response", func(t *testing.T) {
		expectedTitles := map[string]bool{
			"Always":      false,
			"Die For You": false,
		}

		for _, data := range music {
			if _, ok := expectedTitles[data.Title]; ok {
				expectedTitles[data.Title] = true
			}
		}

		for title, found := range expectedTitles {
			if !found {
				t.Errorf("Expected title '%s', not found in response", title)
			}
		}

	})
}

func TestGetMusic(t *testing.T) {
	// Setup test data
	setupTestDB(t)

	// Clear Redis cache
	database.RDB.Del(database.Ctx, "all_musics")

	t.Run("If id exist, return status 200 and its json", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/musics/1", nil)
		recorder := httptest.NewRecorder()
		handler.GetMusic(recorder, req)

		var musicStruct model.Music
		log.Printf("!!!: %v", recorder.Body)
		if err := json.Unmarshal(recorder.Body.Bytes(), &musicStruct); err != nil {
			t.Errorf("Failed to decode response: %v", err)
		}

		got := model.Music{Title: musicStruct.Title, Artist: musicStruct.Artist}
		want := model.Music{Title: "Always", Artist: "Daniel Caesar"}
		if got != want {
			t.Errorf("Expected info: %v, got: %v", want, got)
		}

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status code: %d, got: %d", http.StatusOK, recorder.Code)
		}
	})

	t.Run("If id not exist return status 404", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/musics/0", nil)
		recorder := httptest.NewRecorder()
		handler.GetMusic(recorder, req)

		if recorder.Code != http.StatusNotFound {
			t.Errorf("Expected status code: %d, got: %d", http.StatusNotFound, recorder.Code)
		}
	})
}
