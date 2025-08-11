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

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy from correct path in builder stage
COPY --from=builder /app/standup-logger .

EXPOSE 8080

ENV ENVIRONMENT=production

CMD ["./standup-logger"]
