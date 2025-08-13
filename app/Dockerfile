# ---------- 1) Frontend build (Bun) ----------
FROM oven/bun:1-alpine AS frontend
WORKDIR /app
RUN apk add --no-cache make
COPY . .
RUN make bun-install
RUN make bun-build

# ---------- 2) Go build (+ build migrate) ----------
FROM golang:1.24-alpine AS gobuild
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/standup-logger ./cmd/app/main.go
# build migrate CLI
RUN go install github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0
# (DO NOT set ENTRYPOINT here; this is a build stage only)

# ---------- 3) Runtime ----------
FROM alpine:3.20
WORKDIR /app

# tools you need at runtime
RUN apk add --no-cache ca-certificates make && update-ca-certificates

# app + migrate
COPY --from=gobuild  /app/standup-logger            /app/standup-logger
COPY --from=gobuild  /go/bin/migrate                /usr/local/bin/migrate

# app assets
COPY --from=frontend /app/frontend                  /app/frontend
COPY public                                          /app/public
COPY views                                           /app/views

# migrations + Makefile
COPY db/migrations                                   /app/db/migrations
COPY Makefile                                        /app/Makefile

# entrypoint
COPY entrypoint.sh                                   /entrypoint.sh
RUN chmod +x /entrypoint.sh

# optional switch (default skip)
ENV RUN_MIGRATIONS=false
EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]
