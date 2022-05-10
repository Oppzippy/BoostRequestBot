package database

import (
	"database/sql"
	"time"

	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func (repo *dbRepository) GetBoostRequestChannelByFrontendChannelID(guildID, frontendChannelID string) (*repository.BoostRequestChannel, error) {
	brc, err := repo.getBoostRequestChannel(
		"WHERE guild_id = ? AND frontend_channel_id = ? AND deleted_at IS NULL",
		guildID,
		frontendChannelID,
	)
	return brc, err
}

func (repo *dbRepository) GetBoostRequestChannels(guildID string) ([]*repository.BoostRequestChannel, error) {
	channels, err := repo.getBoostRequestChannels("WHERE guild_id = ? AND frontend_channel_id IS NOT NULL AND deleted_at IS NULL", guildID)
	return channels, err
}

func (repo *dbRepository) getBoostRequestChannel(where string, args ...interface{}) (*repository.BoostRequestChannel, error) {
	channels, err := repo.getBoostRequestChannels(where, args...)
	if err != nil {
		return nil, err
	}
	if len(channels) == 0 {
		return nil, repository.ErrNoResults
	}
	return channels[0], nil
}

func (repo *dbRepository) getBoostRequestChannels(where string, args ...interface{}) ([]*repository.BoostRequestChannel, error) {
	rows, err := repo.db.Query(
		`SELECT id, guild_id, frontend_channel_id, backend_channel_id, uses_buyer_message, skips_buyer_dm
			FROM boost_request_channel `+where,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	channels := make([]*repository.BoostRequestChannel, 0, 1)
	for rows.Next() {
		brc := repository.BoostRequestChannel{}
		var usesBuyerMessage, skipsBuyerDM int
		err := rows.Scan(&brc.ID, &brc.GuildID, &brc.FrontendChannelID, &brc.BackendChannelID, &usesBuyerMessage, &skipsBuyerDM)
		if err != nil {
			return nil, err
		}
		brc.UsesBuyerMessage = usesBuyerMessage != 0
		brc.SkipsBuyerDM = skipsBuyerDM != 0
		channels = append(channels, &brc)
	}

	return channels, nil
}

func (repo *dbRepository) InsertBoostRequestChannel(brc *repository.BoostRequestChannel) error {
	var usesBuyerMessage, skipsBuyerDM int8
	if brc.UsesBuyerMessage {
		usesBuyerMessage = 1
	}
	if brc.SkipsBuyerDM {
		skipsBuyerDM = 1
	}

	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE
			boost_request_channel
		SET
			deleted_at = ?
		WHERE
			guild_id = ? AND
			frontend_channel_id = ? AND
			deleted_at IS NULL`,
		time.Now().UTC(),
		brc.GuildID,
		brc.FrontendChannelID,
	)
	if err = rollbackIfErr(tx, err); err != nil {
		return err
	}

	res, err := tx.Exec(`
		INSERT INTO boost_request_channel (
			guild_id,
			frontend_channel_id,
			backend_channel_id,
			uses_buyer_message,
			skips_buyer_dm,
			created_at
		) VALUES (?, ?, ?, ?, ?, ?)`,
		brc.GuildID,
		sql.NullString{
			String: brc.FrontendChannelID,
			Valid:  brc.FrontendChannelID != "",
		},
		brc.BackendChannelID,
		usesBuyerMessage,
		skipsBuyerDM,
		time.Now().UTC(),
	)
	if err = rollbackIfErr(tx, err); err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err = rollbackIfErr(tx, err); err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	brc.ID = id
	return nil
}

func (repo *dbRepository) DeleteBoostRequestChannel(brc *repository.BoostRequestChannel) error {
	_, err := repo.db.Exec(
		"UPDATE boost_request_channel SET deleted_at = ? WHERE id = ?",
		time.Now().UTC(),
		brc.ID,
	)
	return err
}

func (repo *dbRepository) DeleteBoostRequestChannelsInGuild(guildID string) error {
	_, err := repo.db.Exec(
		"UPDATE boost_request_channel SET deleted_at = ? WHERE guild_id = ?",
		time.Now().UTC(),
		guildID,
	)
	return err
}
