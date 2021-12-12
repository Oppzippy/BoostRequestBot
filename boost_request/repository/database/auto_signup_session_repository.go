package database

import (
	"strings"
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

	return rows.Next(), rows.Err()
}

func (repo *dbRepository) EnableAutoSignup(guildID, advertiserID string, expiresAt time.Time) (*repository.AutoSignUpSession, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
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
		return nil, err
	}
	res, err := tx.Exec(`
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
		return nil, err
	}

	id, err := res.LastInsertId()
	if err := rollbackIfErr(tx, err); err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return &repository.AutoSignUpSession{
		ID: id,
	}, nil
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

func (repo *dbRepository) InsertAutoSignupDelayedMessages(autoSignup *repository.AutoSignUpSession, delayedMessages []*repository.DelayedMessage) error {
	if len(delayedMessages) == 0 {
		return nil
	}
	values := make([]interface{}, 0, len(delayedMessages)*2)
	for _, m := range delayedMessages {
		values = append(values, autoSignup.ID, m.ID)
	}

	_, err := repo.db.Exec(`
		INSERT INTO
			auto_signup_delayed_message (
				auto_signup_id,
				delayed_message_id
			) VALUES (?, ?)`+strings.Repeat(", (?, ?)", len(delayedMessages)-1),
		values...,
	)
	return err
}

func (repo *dbRepository) GetAutoSignupDelayedMessageIDs(guildID string, advertiserID string) ([]int64, error) {
	rows, err := repo.db.Query(`
		SELECT
			asdm.delayed_message_id
		FROM
			auto_signup_delayed_message AS asdm
		INNER JOIN auto_signup_session AS ass
			ON ass.id = asdm.auto_signup_id
		INNER JOIN delayed_message AS dm
			ON dm.id = asdm.delayed_message_id
		WHERE
			ass.guild_id = ? AND
			ass.advertiser_id = ? AND
			dm.deleted_at IS NULL AND
			dm.sent_at IS NULL`,
		guildID,
		advertiserID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]int64, 0, 10)
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
