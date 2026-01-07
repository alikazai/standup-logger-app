package db

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/rs/zerolog/log"
)

func NewDB(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Panic().Err(err).Msg("failed to open DB")
	}

	if err := db.Ping(); err != nil {
		log.Panic().Err(err).Msg("failed to ping DB")
	}

	log.Info().Msg("DB connection established")
	return db
}
