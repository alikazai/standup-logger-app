package db

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func RunSeedFiles(db *sql.DB, seedDir string) error {
	files, err := filepath.Glob(filepath.Join(seedDir, "*.sql"))
	if err != nil {
		return err
	}

	for _, file := range files {
		log.Info().Str("file", file).Msg("seeding")
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		if _, err := db.Exec(string(content)); err != nil {
			return err
		}
	}

	log.Info().Msg("seeding completed successfully")
	return nil
}
