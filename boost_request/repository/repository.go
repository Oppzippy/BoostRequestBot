package repository

import (
	"database/sql"
	"errors"
)

var ErrNoResults = errors.New("not found")
var ErrTooManyResults = errors.New("too many results")

type Repository interface {
	BoostRequestChannelRepository
	BoostRequestRepository
	AdvertiserPrivilegesRepository
	LogChannelRepository
	StealCreditRepository
	ApiKeyRepository
	RoleDiscountRepository
}

type dbRepository struct {
	db *sql.DB
}

func NewDBRepository(db *sql.DB) Repository {
	repo := dbRepository{
		db: db,
	}
	return &repo
}
