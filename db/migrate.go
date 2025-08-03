package db

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/rs/zerolog/log"
)

func RunMigrations(dsn string, migrationsPath string) {
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath), dsn,
	)
	if err != nil {
		log.Panic().Err(err).Msg("failed to initialise migration")
	}

	if err := m.Up(); err != migrate.ErrNoChange {
		log.Panic().Err(err).Msg("migration failed")
	} else {
		log.Info().Msg("database migrated successfully")
	}
}
