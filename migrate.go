package main

import (
	"embed"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/johejo/golang-migrate-extra/source/iofs"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

func MigrateUp(dbURL string) error {
	fs, err := iofs.New(migrationFS, "migrations")
	if err != nil {
		return err
	}
	m, err := migrate.NewWithSourceInstance("iofs", fs, dbURL)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
