package database

import (
	"time"

	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func (repo *dbRepository) IsAutoSignupEnabled(guildID, advertiserID string) (bool, error) {
	rows, err := repo.db.Query(`
		SELECT
			1
		FROM
			auto_signup_session
		WHERE
			guild_id = ? AND
			advertiser_id = ? AND
			expires_at < ? AND
			deleted_at IS NULL`,
		guildID,
		advertiserID,
		time.Now().UTC(),
	)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

func (repo *dbRepository) EnableAutoSignup(guildID, advertiserID string, expiresAt time.Time) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		UPDATE
			auto_signup_session
		SET
			deleted_at = ?
		WHERE
			guild_id = ? AND
			advertiser_id = ? AND
			expires_at > ?`,
		time.Now().UTC(),
		guildID,
		advertiserID,
		time.Now().UTC(),
	)
	if err := rollbackIfErr(tx, err); err != nil {
		return err
	}
	_, err = tx.Exec(`
		INSERT INTO
			auto_signup_session (
				guild_id,
				advertiser_id,
				created_at,
				expires_at
			)
		VALUES (?, ?, ?, ?)`,
		guildID,
		advertiserID,
		time.Now().UTC(),
		expiresAt,
	)
	if err := rollbackIfErr(tx, err); err != nil {
		return err
	}
	return tx.Commit()
}

func (repo *dbRepository) CancelAutoSignup(guildID, advertiserID string) error {
	_, err := repo.db.Exec(`
		UPDATE
			auto_signup_session
		SET
			deleted_at = ?
		WHERE
			guild_id = ? AND
			advertiser_id = ? AND
			expires_at >= ? AND
			deleted_at IS NULL`,
		time.Now().UTC(),
		guildID,
		advertiserID,
		time.Now().UTC(),
	)
	return err
}

func (repo *dbRepository) GetEnabledAutoSignups() ([]*repository.AutoSignupSession, error) {
	sessions, err := repo.getAutoSignups("expires_at >= ? AND deleted_at IS NULL", time.Now().UTC())
	return sessions, err
}

func (repo *dbRepository) GetEnabledAutoSignupsInGuild(guildID string) ([]*repository.AutoSignupSession, error) {
	sessions, err := repo.getAutoSignups(
		"guild_id = ? AND expires_at >= ? AND deleted_at IS NULL",
		guildID,
		time.Now().UTC(),
	)
	return sessions, err
}

func (repo *dbRepository) getAutoSignups(where string, args ...interface{}) ([]*repository.AutoSignupSession, error) {
	query := `
		SELECT
			guild_id,
			advertiser_id,
			MAX(expires_at)
		FROM
			auto_signup_session `
	if where != "" {
		query += " WHERE " + where
	}
	query += " GROUP BY guild_id, advertiser_id"

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := make([]*repository.AutoSignupSession, 0, 50)
	for rows.Next() {
		var session repository.AutoSignupSession
		err := rows.Scan(&session.GuildID, &session.AdvertiserID, &session.ExpiresAt)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, &session)
	}
	return sessions, nil
}
