package database

import (
	"database/sql"
	"strings"

	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

// Returns a slice of preferred advertiser ids or an empty slice if no preferred advertisers are set
func (repo *dbRepository) getPreferredAdvertisers(br *repository.BoostRequest) (map[string]struct{}, error) {
	rows, err := repo.db.Query(
		`SELECT
			discord_user_id
		FROM
			boost_request_preferred_advertiser
		WHERE
			boost_request_id = ?`,
		br.ID,
	)
	if err == sql.ErrNoRows {
		return make(map[string]struct{}), nil
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userIDs := make(map[string]struct{})
	for rows.Next() {
		var userID string
		err := rows.Scan(&userID)
		if err != nil {
			return nil, err
		}
		userIDs[userID] = struct{}{}
	}
	return userIDs, nil
}

func (repo *dbRepository) updatePreferredAdvertisers(tx *sql.Tx, id int64, advertiserIDs map[string]struct{}) error {
	err := repo.deletePreferredAdvertisersExcept(tx, id, advertiserIDs)
	if err != nil {
		return err
	}
	err = repo.insertPreferredAdvertisers(tx, id, advertiserIDs)
	return err
}

func (repo *dbRepository) deletePreferredAdvertisersExcept(tx *sql.Tx, id int64, advertiserIDs map[string]struct{}) error {
	query := "DELETE FROM boost_request_preferred_advertiser WHERE boost_request_id = ?"
	if len(advertiserIDs) > 0 {
		query += " AND discord_user_id NOT IN (?" + strings.Repeat(",?", len(advertiserIDs)-1) + ")"
	}
	args := make([]interface{}, 0, len(advertiserIDs)+1)
	args = append(args, id)
	for advertiserID := range advertiserIDs {
		args = append(args, advertiserID)
	}
	_, err := tx.Exec(query, args...)
	return err
}

func (repo *dbRepository) insertPreferredAdvertisers(tx *sql.Tx, id int64, advertiserIDs map[string]struct{}) error {
	if len(advertiserIDs) == 0 {
		return nil
	}
	query := "INSERT INTO boost_request_preferred_advertiser (boost_request_id, discord_user_id) VALUES (?, ?)"
	query += strings.Repeat(", (?, ?)", len(advertiserIDs)-1)
	args := make([]interface{}, 0, len(advertiserIDs)*2)
	for advertiserID := range advertiserIDs {
		args = append(args, id)
		args = append(args, advertiserID)
	}
	_, err := tx.Exec(query, args...)
	return err
}
