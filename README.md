# Simple Web Server in Go

A simple HTTP server built with Go using only the standary library (`net/http`).

## Features
- GET `/` → returns "Hello World"
- GET `/greet/{name}` → returns "Hi, {name}!"

## How to run
```bash
go run .
```

## How to test
```bash
go test -v
```

---

# Learning points
- Practiced **Test-Driven Development(TDD)** with `httptest`
- Used `http.NewServeMux()` for routing
- Handled dynamic URL path with `strings.TrimPrefix`