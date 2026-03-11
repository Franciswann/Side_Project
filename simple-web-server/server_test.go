package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHelloHandler tests the HelloHandler returns 200 and "Hello World"
func TestHelloHandler(t *testing.T) {
	// create a GET request to "/"
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// create a response recorder to capture the output
	rr := httptest.NewRecorder()

	// call the handler (we haven't written it yet)
	// use type conversion HandlerFunc()
	handler := http.HandlerFunc(HelloHandler)
	handler.ServeHTTP(rr, req)

	// check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// check response body
	expected := "Hello World"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %q want %q", rr.Body.String(), expected)
	}

}

// TestGreetHandler tests that /greet/{name} returns "Hi, {name}!"
func TestGreetHandler(t *testing.T) {
	// create a request to "/greet/Francis"
	req, err := http.NewRequest("GET", "/greet/Francis", nil)
	if err != nil {
		t.Fatal(err)
	}

	// create a response recorder
	rr := httptest.NewRecorder()

	// call the handler (haven't written it yet)
	handler := http.HandlerFunc(GreetHandler)
	handler.ServeHTTP(rr, req)

	// check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// check response body
	expected := "Hi, Francis!"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %q want %q", rr.Body.String(), expected)
	}

}
