package main

import (
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/johejo/golang-migrate-extra/source/iofs"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

func MigrateUp(dbURL string) error {
	fs, err := iofs.New(migrationFS, "migrations")
	if err != nil {
		return fmt.Errorf("opening embedded migration fs: %w", err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", fs, dbURL)
	if err != nil {
		return fmt.Errorf("connecting to db for migrations: %w", err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("running migrations: %w", err)
	}
	return nil
}
