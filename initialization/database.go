package initialization

import (
	"database/sql"
	"fmt"
	"os"
)

func GetDBC() (*sql.DB, error) {
	protocol := os.Getenv("DB_PROTOCOL")
	if protocol == "" {
		protocol = "tcp"
	}

	dataSourceName := fmt.Sprintf(
		"%s:%s@%s(%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		protocol,
		os.Getenv("DB_ADDRESS"),
		os.Getenv("DB_DATABASE"),
	)
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
