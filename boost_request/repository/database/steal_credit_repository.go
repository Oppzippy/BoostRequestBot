package database

import (
	"database/sql"
	"time"

	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func (repo *dbRepository) GetStealCreditsForUser(guildID, userID string) (int, error) {
	row := repo.db.QueryRow(
		`SELECT credits FROM boost_request_steal_credits
		WHERE
			guild_id = ?
			AND user_id = ?
			ORDER BY id DESC
			LIMIT 1`,
		guildID,
		userID,
	)
	var credits int
	err := row.Scan(&credits)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return credits, err
}

func (repo *dbRepository) AdjustStealCreditsForUser(guildID, userID string, operation repository.Operation, amount int) error {
	var operationSymbol string
	switch operation {
	case repository.OperationAdd:
		operationSymbol = "+"
	case repository.OperationSubtract:
		operationSymbol = "-"
	case repository.OperationMultiply:
		operationSymbol = "*"
	case repository.OperationDivide:
		operationSymbol = "/"
	case repository.OperationSet:
		return repo.UpdateStealCreditsForUser(guildID, userID, amount)
	default:
		return repository.ErrInvalidOperation
	}
	_, err := repo.db.Exec(
		`INSERT INTO boost_request_steal_credits (
			guild_id, user_id, credits, created_at
		) VALUES (
			?,
			?,
			COALESCE(
				(SELECT sc.credits FROM boost_request_steal_credits sc WHERE sc.guild_id = ? AND sc.user_id = ? ORDER BY sc.id DESC LIMIT 1),
				0
			) `+operationSymbol+` ?,
			?
		)`,
		guildID,
		userID,
		guildID,
		userID,
		amount,
		time.Now().UTC(),
	)
	return err
}

func (repo *dbRepository) UpdateStealCreditsForUser(guildID, userID string, amount int) error {
	_, err := repo.db.Exec(
		`INSERT INTO boost_request_steal_credits (guild_id, user_id, credits, created_at) VALUES (?, ?, ?, ?)`,
		guildID,
		userID,
		amount,
		time.Now().UTC(),
	)
	return err
}
