package database

import (
	"database/sql"
)

type dbRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *dbRepository {
	repo := dbRepository{
		db: db,
	}
	return &repo
}
