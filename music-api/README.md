# Music API - Golang RESTful Backend

A simple REST API for managing music tracks, built with Go standard library only (no external frameworks).

## Features
- GET `/musics` - List all musics
- POST `/musics` - Create a new music
- GET `/musics/{id}` - Get a specific music
- DELETE `/musics/{id}` - Delete a music
- PUT `/musics/{id}` - Update a music

## Tech Stack
- Go (standard library)
- `net/http` + `encoding/json`
- In-memory storage with `map`
- Tdd with `httptest`

## How to Run
``` bash
go run .
```

## How to Test
``` bash
go test -v
```

## API Examples
``` bash
# Get all musics
curl http://localhost:8080/musics

# Create a new music
curl -X POST -H "Content-Type: application/json" \
    -d '{"title":"Blinding Lights","artist":"The Weekend"}' \
    http://localhost:8080/musics
```

## Learning Points
- Built full CRUD operations using only stdlib
- Practiced TDD with `httptest`
- Handled JSON encoding/decoding and proper HTTP status codes
- Managed in-memory data with map and ID generation