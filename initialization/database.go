package initialization

import (
	"database/sql"
	"fmt"
	"os"
)

func GetDBC() (*sql.DB, error) {
	dataSourceName := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)
	if dsn := os.Getenv("DB_DSN"); dsn != "" {
		dataSourceName = dsn
	}

	db, err := sql.Open("mysql", dataSourceName+"?parseTime=true")
	if err != nil {
		return nil, fmt.Errorf("connecting to db: %v", err)
	}

	err = MigrateUp("mysql://" + dataSourceName + "?multiStatements=true")
	if err != nil {
		return nil, fmt.Errorf("running migrations: %v", err)
	}
	return db, nil
}
