package repository

import (
	"database/sql"
)

type Repository interface {
	BoostRequestChannelRepository
	BoostRequestRepository
	AdvertiserPrivilegesRepository
	LogChannelRepository
	StealCreditRepository
}

type dbRepository struct {
	db *sql.DB
}

func NewDBRepository(db *sql.DB) Repository {
	repo := dbRepository{
		db: db,
	}
	return repo
}
