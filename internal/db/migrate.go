package db

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate"
)

func MigrateDatabase(dbURL, direction string) error {

	m, err := migrate.New(
		"file://../internal/db/migrations",
		dbURL,
	)
	if err != nil {
		log.Fatal("creating migrate instance: %w", err)
	}

	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("migrating up: %w", err)
		}
	case "down":
		if err := m.Down(); err != nil {
			return fmt.Errorf("migrating down: %w", err)
		}
	case "force":
		version, _, err := m.Version()
		if err != nil {
			return fmt.Errorf("getting current version: %w", err)
		}

		if err := m.Force(int(version)); err != nil {
			return fmt.Errorf("forcing version: %w", err)
		}
	default:
		return fmt.Errorf("invalid migration direction: %s", direction)
	}

	return nil
}
