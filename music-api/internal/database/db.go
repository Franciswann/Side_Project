package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"music-api/internal/model"

	"github.com/redis/go-redis/v9"
)

var (
	DB  *sql.DB
	RDB *redis.Client
	Ctx = context.Background()
)

var musics = make(map[int]model.Music)

func init() {
	musics[1] = model.Music{
		Id:     1,
		Title:  "Perfect",
		Artist: "Ed Sheeran",
	}
	musics[2] = model.Music{
		Id:     2,
		Title:  "Always",
		Artist: "Daniel Caesar",
	}
	musics[3] = model.Music{
		Id:     3,
		Title:  "Die For You",
		Artist: "Joji",
	}
}

// Initialize database setup
func InitDB(connStr string) error {
	// connecting to pq driver
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	// defer DB.Close()

	// Test connection
	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connected to PostgreSQL successfully")
	}

	// insert initial data to table 'musics' with conflict handling
	for _, music := range musics {
		query := `INSERT INTO musics (title, artist) 
				  VALUES ($1, $2) 
				  ON CONFLICT (title, artist) DO NOTHING;`
		_, err := DB.Exec(query, music.Title, music.Artist)
		if err != nil {
			return err
		}
	}
	log.Println("Insert initial data successfully")
	return nil
}

// Initialize Redis connection
// Use Ping() to test if Redis is connected
func InitRedis(redisHost, redisPort string) error {

	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: "",
		DB:       0,
	})

	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		log.Printf("Failed to connect to Redis: %s", err)
	} else {
		log.Printf("Connected to Redis successfully")
	}
	return nil
}
