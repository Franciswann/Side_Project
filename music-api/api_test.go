package main

import (
	"encoding/json"
	"maps"
	"net/http"
	"net/http/httptest"
	"strings"
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

	var newMusic Music

	// back up global map 'musics'
	originalMusics := maps.Clone(musics)

	jsonBody := `{"title":"Sparks","artist":"Coldplay"}`

	// request contains a io.Reader from strings.NewReader
	req := httptest.NewRequest("POST", "/musics", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	CreateMusic(recorder, req)

	// decode the response body to get the created music with ID
	if err := json.NewDecoder(recorder.Body).Decode(&newMusic); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	t.Run("check Header's validation", func(t *testing.T) {
		if got := recorder.Header().Get("Content-Type"); got != "application/json" {
			t.Errorf("got %s, want %s", got, "application/json")
		}
	})

	t.Run("test status code = StatusCreated 201", func(t *testing.T) {
		if got := recorder.Code; got != http.StatusCreated {
			t.Errorf("got %v, want %v", got, http.StatusCreated)
		}
	})

	t.Run("check if the music ID is the newest one", func(t *testing.T) {
		if newMusic.Id != len(originalMusics)+1 {
			t.Errorf("got %d, want %d", newMusic.Id, len(musics))
		}
	})

	//  (using parsed response, not map)
	t.Run("verify the created music matches expectations", func(t *testing.T) {
		want := Music{
			Id:     len(originalMusics) + 1,
			Title:  "Sparks",
			Artist: "Coldplay",
		}
		if newMusic != want {
			t.Errorf("got %v, want %v", newMusic, want)
		}
	})

	// restore data
	musics = maps.Clone(originalMusics)

}
