# syntax=docker/dockerfile:1

# Stage 1: Build and compile
# Use the golang image from Docker Hub
FROM golang:1.25.6-alpine AS base
# Instruct Docker to use this DIR as the default destination
WORKDIR /app

# Copy dependencies to WORKDIR
COPY go.mod go.sum ./

# Install app dependencies
RUN go mod download

# Copy the rest of the app into the container
COPY . .

# Set port environment variable
ENV PORT=8080
# Open port for host
EXPOSE 8080

# RUN run the command in container

# Build a static binary file (main) from cmd/app/main.go and output
# to /app/music-api and then run it in a minimal final stage (alpine)
RUN go build -o /music-api ./cmd/app/

# ========================================

# Stage 2: Run
# Use Alpine Linux as environment
FROM alpine:latest

# Run command on /app
WORKDIR /app
# Cpoy binary
COPY --from=base /music-api .

CMD ["./music-api"]