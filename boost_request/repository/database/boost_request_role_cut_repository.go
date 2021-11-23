package database

import (
	"database/sql"
	"strings"
)

func (repo *dbRepository) getRoleCuts(boostRequestID int64) (map[string]int64, error) {
	rows, err := repo.db.Query(`
		SELECT
			role_id,
			role_cut
		FROM
			boost_request_role_cut
		WHERE
			boost_request_id = ?`,
		boostRequestID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cutsByRoleID := make(map[string]int64)
	for rows.Next() {
		var roleID string
		var cut int64
		err := rows.Scan(&roleID, &cut)
		if err != nil {
			return nil, err
		}
		cutsByRoleID[roleID] = cut
	}
	return cutsByRoleID, nil
}

func (repo *dbRepository) updateRoleCuts(tx *sql.Tx, boostRequestID int64, cuts map[string]int64) error {
	err := repo.deleteRoleCutsExcept(tx, boostRequestID, cuts)
	if err != nil {
		return err
	}
	err = repo.insertRoleCuts(tx, boostRequestID, cuts)
	return err
}

func (repo *dbRepository) deleteRoleCutsExcept(tx *sql.Tx, boostRequestID int64, cuts map[string]int64) error {
	args := make([]interface{}, 0, len(cuts)+1)
	args = append(args, boostRequestID)
	for roleID := range cuts {
		args = append(args, roleID)
	}
	query := "DELETE FROM boost_request_role_cut WHERE boost_request_id = ?"
	if len(cuts) > 0 {
		query += " AND role_id NOT IN (?" + strings.Repeat(",?", len(cuts)-1) + ")"
	}
	_, err := tx.Exec(query, args...)
	return err
}

func (repo *dbRepository) insertRoleCuts(tx *sql.Tx, boostRequestID int64, cuts map[string]int64) error {
	if len(cuts) == 0 {
		return nil
	}
	query := "INSERT INTO boost_request_role_cut (boost_request_id, role_id, role_cut) VALUES (?, ?, ?)"
	query += strings.Repeat(", (?, ?, ?)", len(cuts)-1)
	args := make([]interface{}, 0, len(cuts)*3)
	for roleID, cut := range cuts {
		args = append(args, boostRequestID, roleID, cut)
	}
	_, err := tx.Exec(query, args...)
	return err
}
