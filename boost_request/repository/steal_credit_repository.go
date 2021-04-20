package repository

import (
	"database/sql"
	"errors"
	"time"
)

type StealCreditRepository interface {
	GetStealCreditsForUser(guildID, userID string) (int, error)
	AdjustStealCreditsForUser(guildID, userID string, operation Operation, amount int) error
	UpdateStealCreditsForUser(guildID, userID string, amount int) error
}

var ErrInvalidOperation = errors.New("invalid math operation")

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

func (repo *dbRepository) AdjustStealCreditsForUser(guildID, userID string, operation Operation, amount int) error {
	var operationSymbol string
	switch operation {
	case OperationAdd:
		operationSymbol = "+"
	case OperationSubtract:
		operationSymbol = "-"
	case OperationMultiply:
		operationSymbol = "*"
	case OperationDivide:
		operationSymbol = "/"
	case OperationSet:
		return repo.UpdateStealCreditsForUser(guildID, userID, amount)
	default:
		return ErrInvalidOperation
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
