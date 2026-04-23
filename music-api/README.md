# Music API v2 - Production Ready Golang Backend

A RESTful Music API built with Go, upgraded from in-memory storage to PostgreSQL + Redis cache, and containerized with Docker.

## Features
- GET `/musics` - List all musics with Redis caching
- POST `/musics` - Create a new music
- GET `/musics/{id}` - Get a specific music with Redis caching
- DELETE `/musics/{id}` - Delete a music
- PUT `/musics/{id}` - Update a music

## Tech Stack
- **Go**: Standard library only, no external frameworks
- **PostgreSQL**: Persistent storage with connection pooling
- **Redis**: Caching layer with automatic invalidation
- **Testing**: Integration tests with TestMain setup for list & get endpoints
- **Docker**: Multi-stage builds and orchestration
- **API**: RESTful endpoints with proper error handling

## How to Run

**Local Development**
``` bash
go run .
```
or

**Docker Development**
``` bash
docker compose up --build
```
![alt text](image-5.png)

![alt text](image-3.png)

### Stop Containers
Clean shutdown of all services (API, PostgreSQL, Redis).


``` bash
docker compose down
```
![alt text](image.png)

## How to Test
**Testing**
``` bash
go test -v
```

## API Examples
``` bash
# Get all musics
curl -v -X GET http://localhost:8080/musics

# Get specified music
curl -v -X GET http://localhost:8080/musics/2

# Create a new music
curl -v -X POST -H "Content-Type: application/json" \
    -d '{"title":"Blinding Lights","artist":"The Weekend"}' \
    http://localhost:8080/musics

# Delete a music
curl -v -X DELETE http://localhost:8080/musics/3
```

Example of `GET`:
![alt text](image-8.png)

## Learning Points
- Built full **CRUD** operations with **PostgreSQL**
- Implemented **Redis cache** with invalidation
- Handled **JSON** encoding/decoding and proper **HTTP status codes**
- Practiced **TDD** and wrote integration tests with `httptest`
- Containerized the application with **Docker**