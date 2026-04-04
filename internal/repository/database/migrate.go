package database

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"

	myErrors "TaskTracker/errors"
)

func runMigrations(connStr string) error {
	migration, err := migrate.New(
		"file://../migrations",
		connStr,
	)

	if err != nil {
		return fmt.Errorf("%v : %w", myErrors.ErrCreateMigration, err)
	}

	defer migration.Close()

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("%v : %w", myErrors.ErrCantUseMigration, err)
	}
	return nil
}
