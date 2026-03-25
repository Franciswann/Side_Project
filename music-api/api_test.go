package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListMusics(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/musics", nil)
	recorder := httptest.NewRecorder()

	// because HandlerFunc already implement ServeHTTP, so
	// ListMusics can seen as Handler
	ListMusics(recorder, req)

	t.Run("check http status code", func(t *testing.T) {

		if status := recorder.Code; status != http.StatusOK {
			t.Errorf("got %v, want %v", status, http.StatusOK)
		}
	})

	t.Run("check the data has been writen", func(t *testing.T) {
		var music []Music
		err := json.Unmarshal(recorder.Body.Bytes(), &music)
		if err != nil {
			t.Errorf("%v", err)
		}

		got := len(music)
		want := 3
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
}

// create an body with header and then test its result, status code
func TestCreateMusic(t *testing.T) {
	var music = Music{
		Title:  "Sparks",
		Artist: "Coldplay",
	}

	data, err := json.Marshal(music)
	if err != nil {
		t.Errorf("Marshaling failed")
	}
	// request contains a io.Reader from strings.NewReader
	req := httptest.NewRequest("POST", "/musics", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	CreateMusic(recorder, req)

	// decode the response body to get the created music with ID
	if err := json.NewDecoder(recorder.Body).Decode(&music); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	t.Run("test status code = StatusCreated 201", func(t *testing.T) {
		if got := recorder.Code; got != http.StatusCreated {
			t.Errorf("got %v, want %v", got, http.StatusCreated)
		}
	})

	t.Run("check if the music ID is the newest one", func(t *testing.T) {
		if music.Id != 4 {
			t.Errorf("got %d, want %d", music.Id, 4)
		}
	})

	//  (using parsed response, not map)
	t.Run("verify the created music matches expectations", func(t *testing.T) {
		want := Music{
			Id:     music.Id,
			Title:  "Sparks",
			Artist: "Coldplay",
		}
		if music != want {
			t.Errorf("got %v, want %v", music, want)
		}
	})
}

// curl -v -X POST -H "Content-Type: application/json" \
// -d '{"title":"Blinding Lights","artist":"The Weekend"}' \
// http://localhost:8080/musics
