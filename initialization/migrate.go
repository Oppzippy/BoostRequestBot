package initialization

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/oppzippy/BoostRequestBot/migrations"
)

func MigrateUp(dbURL string) error {
	fs, err := iofs.New(migrations.MigrationFS, ".")
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
