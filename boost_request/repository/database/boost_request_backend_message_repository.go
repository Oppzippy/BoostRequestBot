package database

import (
	"database/sql"
	"strings"

	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func (repo *dbRepository) getBoostRequestBackendMessages(boostRequestID int64) ([]*repository.BoostRequestBackendMessage, error) {
	rows, err := repo.db.Query(`
		SELECT
			channel_id,
			message_id
		FROM
			boost_request_backend_message
		WHERE
			boost_request_id = ?`,
		boostRequestID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]*repository.BoostRequestBackendMessage, 0, 1)
	for rows.Next() {
		m := repository.BoostRequestBackendMessage{}
		err := rows.Scan(&m.ChannelID, &m.MessageID)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &m)
	}
	return messages, nil
}

func (repo *dbRepository) updateBoostRequestBackendMessages(tx *sql.Tx, boostRequestID int64, messages []*repository.BoostRequestBackendMessage) error {
	err := repo.deleteBoostRequestBackendMessages(tx, boostRequestID)
	if err != nil {
		return err
	}
	err = repo.insertBoostRequestBackendMessages(tx, boostRequestID, messages)
	return err
}

func (repo *dbRepository) deleteBoostRequestBackendMessages(tx *sql.Tx, boostRequestID int64) error {
	_, err := tx.Exec(
		`DELETE FROM boost_request_backend_message WHERE boost_request_id = ?`,
		boostRequestID,
	)
	return err
}
func (repo *dbRepository) insertBoostRequestBackendMessages(tx *sql.Tx, boostRequestID int64, messages []*repository.BoostRequestBackendMessage) error {
	if len(messages) > 0 {
		query := `INSERT INTO
				boost_request_backend_message (
					boost_request_id,
					channel_id,
					message_id
				) VALUES (?, ?, ?) ` + strings.Repeat(", (?, ?, ?)", len(messages)-1)
		args := make([]interface{}, 0, len(messages)*3)
		for _, m := range messages {
			args = append(args, boostRequestID, m.ChannelID, m.MessageID)
		}
		_, err := tx.Exec(query, args...)
		if err != nil {
			return err
		}
	}
	return nil
}
