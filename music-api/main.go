package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

// This Music API can let user Create, Read, Update nad Delete musics
func main() {

	// connecting to pq driver
	connStr := "host=localhost port=5432 user=wanchaochun dbname=music_db sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connected to PostgreSQL")
	}

	err = initDB()
	if err != nil {
		log.Println("insert failed:", err)
	} else {
		fmt.Println("insert succeed")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /musics", ListMusics)
	// GetMusic handles GET requests to /musics/{id}
	// Example: GET /musics/4
	mux.HandleFunc("GET /musics/{id}", GetMusic)
	mux.HandleFunc("POST /musics", CreateMusic)
	mux.HandleFunc("DELETE /musics/{id}", DeleteMusic)
	mux.HandleFunc("PUT /musics/{id}", UpdateMusic)

	log.Println("Running...")

	// ListenAndServe uses the configurable mux
	// test it on terminal: curl http://localhost:8080/
	http.ListenAndServe(":8080", mux)

}
