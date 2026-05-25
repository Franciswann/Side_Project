# Side Projects

This repository contains my side projects for practicing Go and backend development.

## Projects

### 🚀 Production-Ready Projects

- **[music-api](./music-api)**: A production-ready REST API for managing music tracks
    - **Tech Stack**: Go, PostgreSQL, Redis, Docker, Docker Compose
    - **Features**: Complete CRUD operations, caching, containerization, testing
    - **Architecture**: Industry-standard project structure with layered design
    - **Database**: Basic table design with indexes and constraints

- **[task-worker](./task-worker)**: A concurrent job queue system with retry and graceful shutdown
    - **Tech Stack**: Go, `sync`, `context`, `os/signal`
    - **Features**: Worker pool, retry mechanism, dead-letter queue, graceful shutdown
    - **Architecture**: Producer/consumer pattern with dependency injection for testability
    - **Concurrency**: Goroutines, buffered channels, WaitGroup, race-condition-free

### 📚 Learning Projects

- **[simple-web-server](./simple-web-server)**: A simple HTTP server with JSON API and TDD.

## 🎯 Focus Areas

- **Backend Development**: REST APIs, job queues, caching strategies
- **Concurrency**: Goroutines, channels, worker pools, graceful shutdown
- **DevOps**: Docker containerization, service orchestration
- **Testing**: Unit tests, TDD practices, race detector
- **Architecture**: Clean code, layered design, dependency injection


## 📊 Project Statistics
 
- **Total Projects**: 3
- **Production Ready**: 2
- **Languages**: Go
- **Technologies**: PostgreSQL, Redis, Docker, Docker Compose, `sync`, `context`