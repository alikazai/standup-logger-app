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
# build migrate CLI with Postgres driver
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0

# ---------- 3) Runtime ----------
FROM alpine:latest
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

# default: don't run migrations automatically
ENV RUN_MIGRATIONS=false
# optional: when true, run migrations then EXIT (useful for one-off migration job)
ENV RUN_MIGRATIONS_EXIT=false
EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]
