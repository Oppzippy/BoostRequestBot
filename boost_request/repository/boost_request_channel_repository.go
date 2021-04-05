package repository

import (
	"database/sql"
	"errors"
	"time"
)

type BoostRequestChannelRepository interface {
	GetBoostRequestChannelByFrontendChannelID(guildID string, frontendChannelID string) (*BoostRequestChannel, error)
	InsertBoostRequestChannel(brc *BoostRequestChannel) error
}

var ErrBoostRequestChannelNotFound = errors.New("boost request channel not found")

func (repo dbRepository) GetBoostRequestChannelByFrontendChannelID(guildID string, frontendChannelID string) (*BoostRequestChannel, error) {
	brc, err := repo.getBoostRequestChannel(
		"WHERE guild_id = ? AND frontend_channel_id = ? AND deleted_at IS NULL",
		guildID,
		frontendChannelID,
	)
	return brc, err
}

func (repo dbRepository) getBoostRequestChannel(where string, args ...interface{}) (*BoostRequestChannel, error) {
	row := repo.db.QueryRow(
		`SELECT id, guild_id, frontend_channel_id, backend_channel_id, uses_buyer_message, skips_buyer_dm
			FROM boost_request_channel `+where,
		args...,
	)
	brc := BoostRequestChannel{}
	var usesBuyerMessage, skipsBuyerDM int
	err := row.Scan(&brc.ID, &brc.GuildID, &brc.FrontendChannelID, &brc.BackendChannelID, &usesBuyerMessage, &skipsBuyerDM)
	if err == sql.ErrNoRows {
		return nil, ErrBoostRequestChannelNotFound
	}
	if err != nil {
		return nil, err
	}
	brc.UsesBuyerMessage = usesBuyerMessage != 0
	brc.SkipsBuyerDM = skipsBuyerDM != 0

	return &brc, nil
}

func (repo dbRepository) InsertBoostRequestChannel(brc *BoostRequestChannel) error {
	var usesBuyerMessage, skipsBuyerDM int8
	if brc.UsesBuyerMessage {
		usesBuyerMessage = 1
	}
	if brc.SkipsBuyerDM {
		skipsBuyerDM = 1
	}

	res, err := repo.db.Exec(
		`INSERT INTO boost_request_channel (
			guild_id,
			frontend_channel_id,
			backend_channel_id,
			uses_buyer_message,
			skips_buyer_dm,
			created_at
		) VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			backend_channel_id = VALUES(backend_channel_id),
			uses_buyer_message = VALUES(uses_buyer_message),
			skips_buyer_dm = VALUES(skips_buyer_dm)`,
		brc.GuildID,
		brc.FrontendChannelID,
		brc.BackendChannelID,
		usesBuyerMessage,
		skipsBuyerDM,
		time.Now().UTC().Format(time.RFC3339),
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	brc.ID = id
	return nil
}
