# syntax=docker/dockerfile:1

# Stage 1: Build and compile
# Use the golang image from Docker Hub
FROM golang:1.25-alpine AS base
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

# build a static binary file (main) into root(/) and then
# run it in a minimal final stage (alpine)
RUN go build -o /music-api .

# ========================================

# Stage 2: Run
# Use Alpine Linux as environment
FROM alpine:latest

# Run command on /root/
WORKDIR /root/

COPY --from=base /music-api .

CMD ["./music-api"]