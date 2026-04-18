package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var db *sql.DB

// Global Redis Client
var rdb *redis.Client

// Initialize Redis connection
// Use Ping() to test if Redis is connected
func initRedis() {
	// Retrieve the value from .env
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("connection failed: %s", err)
	} else {
		log.Printf("Connected to Redis successfully")
	}
}

// This Music API can let user Create, Read, Update nad Delete musics
func main() {

	// Load environment variables from .env.local for local development
	loadErr := godotenv.Load(".env.local")
	if loadErr != nil {
		log.Printf("Warning: .env.local not found, using default environment variables")
	}

	// Initialize Redis
	initRedis()

	// Initialize PostgreSQL
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASSWORD")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)

	// connecting to pq driver
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connected to PostgreSQL successfully")
	}

	// Initialize database
	err = initDB()
	if err != nil {
		log.Println("insert failed:", err)
	} else {
		log.Println("insert initial data successfully")
	}

	// route setting
	mux := http.NewServeMux()
	mux.HandleFunc("GET /musics", ListMusics)
	// GetMusic handles GET requests to /musics/{id}
	// Example: GET /musics/4
	mux.HandleFunc("GET /musics/{id}", GetMusic)
	mux.HandleFunc("POST /musics", CreateMusic)
	mux.HandleFunc("DELETE /musics/{id}", DeleteMusic)
	mux.HandleFunc("PUT /musics/{id}", UpdateMusic)

	log.Println("Server is running on: 8080...")

	// ListenAndServe uses the configurable mux
	// test it on terminal: curl http://localhost:8080/
	http.ListenAndServe(":8080", mux)

}
