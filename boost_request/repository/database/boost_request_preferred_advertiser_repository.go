package database

import (
	"database/sql"
	"strings"

	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func (repo *dbRepository) GetPreferredAdvertisers(br *repository.BoostRequest) ([]string, error) {
	rows, err := repo.db.Query(
		`SELECT
			discord_user_id
		FROM
			boost_request_preferred_advertiser
		WHERE
			boost_request_id = ?`,
		br.ID,
	)
	if err != nil {
		return nil, err
	}

	userIDs := make([]string, 0, 1)
	for rows.Next() {
		var userID string
		err := rows.Scan(&userID)
		if err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}
	return userIDs, nil
}

func (repo *dbRepository) SetPreferredAdvertisers(br *repository.BoostRequest, advertiserIDs []string) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	err = repo.deletePreferredAdvertisersExcept(tx, br, advertiserIDs)
	if err = rollbackIfErr(tx, err); err != nil {
		return err
	}
	err = repo.insertPreferredAdvertisers(tx, br, advertiserIDs)
	if err = rollbackIfErr(tx, err); err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

func (repo *dbRepository) deletePreferredAdvertisersExcept(tx *sql.Tx, br *repository.BoostRequest, advertiserIDs []string) error {
	query := "DELETE FROM boost_request_preferred_advertiser WHERE boost_request_id = ?"
	if len(advertiserIDs) > 0 {
		query += " AND NOT IN (?" + strings.Repeat(",?", len(advertiserIDs)-1) + ")"
	}
	args := make([]interface{}, 0, len(advertiserIDs)+1)
	args = append(args, br.ID)
	for _, advertiserID := range advertiserIDs {
		args = append(args, advertiserID)
	}
	_, err := tx.Exec(query, args...)
	return err
}

func (repo *dbRepository) insertPreferredAdvertisers(tx *sql.Tx, br *repository.BoostRequest, advertiserIDs []string) error {
	if len(advertiserIDs) == 0 {
		return nil
	}
	query := "INSERT INTO boost_request_preferred_advertiser (boost_request_id, discord_user_id) VALUES (?, ?)"
	query += strings.Repeat(", (?, ?)", len(advertiserIDs)-1)
	args := make([]interface{}, 0, len(advertiserIDs)*2)
	for _, advertiserID := range advertiserIDs {
		args = append(args, br.ID)
		args = append(args, advertiserID)
	}
	_, err := tx.Exec(query, args...)
	return err
}
