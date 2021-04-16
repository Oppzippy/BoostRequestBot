package repository

import (
	"database/sql"
	"errors"
	"time"
)

type BoostRequestChannelRepository interface {
	GetBoostRequestChannelByFrontendChannelID(guildID string, frontendChannelID string) (*BoostRequestChannel, error)
	InsertBoostRequestChannel(brc *BoostRequestChannel) error
	DeleteBoostRequestChannel(brc *BoostRequestChannel) error
	DeleteBoostRequestChannelsInGuild(guildID string) error
	GetBoostRequestChannels(id string) ([]*BoostRequestChannel, error)
}

var ErrBoostRequestChannelNotFound = errors.New("boost request channel not found")

func (repo *dbRepository) GetBoostRequestChannelByFrontendChannelID(guildID, frontendChannelID string) (*BoostRequestChannel, error) {
	brc, err := repo.getBoostRequestChannel(
		"WHERE guild_id = ? AND frontend_channel_id = ?",
		guildID,
		frontendChannelID,
	)
	return brc, err
}

func (repo *dbRepository) GetBoostRequestChannels(guildID string) ([]*BoostRequestChannel, error) {
	channels, err := repo.getBoostRequestChannels("WHERE guild_id = ?", guildID)
	return channels, err
}

func (repo *dbRepository) getBoostRequestChannel(where string, args ...interface{}) (*BoostRequestChannel, error) {
	channels, err := repo.getBoostRequestChannels(where, args...)
	if err != nil {
		return nil, err
	}
	switch len(channels) {
	case 0:
		return nil, ErrBoostRequestChannelNotFound
	case 1:
		return channels[0], nil
	default:
		return nil, ErrTooManyResults
	}
}

func (repo *dbRepository) getBoostRequestChannels(where string, args ...interface{}) ([]*BoostRequestChannel, error) {
	rows, err := repo.db.Query(
		`SELECT id, guild_id, frontend_channel_id, backend_channel_id, uses_buyer_message, skips_buyer_dm
			FROM boost_request_channel `+where,
		args...,
	)
	if err != nil {
		return nil, err
	}
	channels := make([]*BoostRequestChannel, 0, 1)
	for rows.Next() {
		brc := BoostRequestChannel{}
		var usesBuyerMessage, skipsBuyerDM int
		err := rows.Scan(&brc.ID, &brc.GuildID, &brc.FrontendChannelID, &brc.BackendChannelID, &usesBuyerMessage, &skipsBuyerDM)
		if err == sql.ErrNoRows {
			return nil, ErrBoostRequestChannelNotFound
		}
		if err != nil {
			return nil, err
		}
		brc.UsesBuyerMessage = usesBuyerMessage != 0
		brc.SkipsBuyerDM = skipsBuyerDM != 0
		channels = append(channels, &brc)
	}

	return channels, nil
}

func (repo *dbRepository) InsertBoostRequestChannel(brc *BoostRequestChannel) error {
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
		time.Now().UTC(),
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

func (repo *dbRepository) DeleteBoostRequestChannel(brc *BoostRequestChannel) error {
	_, err := repo.db.Exec("DELETE FROM boost_request_channel WHERE id = ?", brc.ID)
	return err
}

func (repo *dbRepository) DeleteBoostRequestChannelsInGuild(guildID string) error {
	_, err := repo.db.Exec("DELETE FROM boost_request_channel WHERE guild_id = ?", guildID)
	return err
}
