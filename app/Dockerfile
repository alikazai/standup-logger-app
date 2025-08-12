# ---------- 1) Frontend build (Bun) ----------
FROM oven/bun:1-alpine AS frontend
WORKDIR /app
# If your Makefile target runs "cd frontend && bun run build",
# copy the whole repo and run make:
RUN apk add --no-cache make
COPY . .
RUN make bun-install
RUN make bun-build


# ---------- 2) Go build (+ build migrate) ----------
FROM golang:1.24-alpine AS gobuild
WORKDIR /app
RUN apk add --no-cache git

# cache modules
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .

# build your app
RUN go build -o /app/standup-logger ./cmd/app/main.go

# build the migrate CLI from source (no GitHub releases download)
# pin the version you want
RUN go install github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0


# ---------- 3) Runtime ----------
FROM alpine:latest
WORKDIR /app

# minimal tools (keep make if you use `make migrate` in ECS one-offs)
RUN apk add --no-cache ca-certificates make && update-ca-certificates

# app binary
COPY --from=gobuild  /app/standup-logger .

# migrate binary (built in gobuild stage)
COPY --from=gobuild  /go/bin/migrate /usr/local/bin/migrate

# migrations + Makefile for one-off ECS task (make migrate)
COPY db/migrations   /app/db/migrations
COPY Makefile        /app/Makefile

# static assets built by Bun (adjust paths to your app’s expectations)
COPY --from=frontend /app/frontend /app/frontend
COPY public /app/public

# copy templates if you serve them
COPY views /app/views

EXPOSE 8080
ENV ENVIRONMENT=production

# default: run web server (override for migrations)
CMD ["./standup-logger"]
