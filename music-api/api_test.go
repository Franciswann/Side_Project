package main

import (
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
			t.Errorf("something goes wrong")
		}

		got := len(music)
		want := 3
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
}

// func TestCreateMusic(t *testing.T) {

// }

// 寫 TestCreateMusic：測試 POST /musics 是否能成功新增並回 201
