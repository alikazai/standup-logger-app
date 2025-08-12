# Dockerfile

# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Copy go.mod and go.sum first to leverage caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the full source
COPY . .

# Build the binary in current dir (not inside /app/app)
RUN go build -o standup-logger ./cmd/app/main.go

# Run stage
FROM alpine:latest

# tools needed at runtime for migrations
RUN apk add --no-cache ca-certificates curl make \
 && curl -sSL https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz \
 |  tar -xz -C /usr/local/bin

WORKDIR /app

# Copy from correct path in builder stage
COPY --from=builder /app/standup-logger .

# copy migrations + Makefile so `make migrate` can run inside container
COPY db/migrations /app/db/migrations
COPY Makefile /app/Makefile

EXPOSE 8080

ENV ENVIRONMENT=production

# default: run web server (override for migrations)
CMD ["./standup-logger"]
