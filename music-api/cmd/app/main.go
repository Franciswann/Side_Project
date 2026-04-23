package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"music-api/internal/database"
	"music-api/internal/handler"

	"github.com/joho/godotenv"
)

// This Music API can let user Create, Read, Update nad Delete musics
func main() {
	// Load environment variables from .env.local for local development
	loadErr := godotenv.Load(".env.local")
	if loadErr != nil {
		log.Printf("Warning: .env.local not found, using default environment variables")
	}

	// Retrieve the value from .env
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASSWORD")
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)

	// Initialize PostgreSQL
	if err := database.InitDB(connStr); err != nil {
		log.Fatal(err)
	}

	// Initialize Redis
	if err := database.InitRedis(redisHost, redisPort); err != nil {
		log.Fatal(err)
	}

	// route setting
	mux := http.NewServeMux()
	mux.HandleFunc("GET /musics", handler.ListMusics)
	// GetMusic handles GET requests to /musics/{id}
	// Example: GET /musics/4
	mux.HandleFunc("GET /musics/{id}", handler.GetMusic)
	mux.HandleFunc("POST /musics", handler.CreateMusic)
	mux.HandleFunc("DELETE /musics/{id}", handler.DeleteMusic)
	mux.HandleFunc("PUT /musics/{id}", handler.UpdateMusic)

	log.Println("Server is running on: 8080...")

	// ListenAndServe uses the configurable mux
	// test it on terminal: curl http://localhost:8080/
	http.ListenAndServe(":8080", mux)

}
