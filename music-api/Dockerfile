# syntax=docker/dockerfile:1
# Use the golang image from Docker Hub
FROM golang:1.25-alpine AS base
# Instruct Docker to use this DIR as the default destination
WORKDIR /app

# Copy dependencies to WORKDIR and download it
COPY go.mod go.sum ./
RUN go mod download

# Copy source code into WORKDIR
COPY *.go ./

# build a static application binary file
RUN go build -o /music-api

EXPOSE 8080

CMD ["/music-api"]