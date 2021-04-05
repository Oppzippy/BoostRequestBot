package repository

import (
	"database/sql"
	"errors"
)

type BoostRequestChannelRepository interface {
	GetBoostRequestChannelByFrontendChannelID(guildID string, frontendChannelID string) (*BoostRequestChannel, error)
}

var ErrBoostRequestChannelNotFound = errors.New("boost request channel not found")

func (repo dbRepository) GetBoostRequestChannelByFrontendChannelID(guildID string, frontendChannelID string) (*BoostRequestChannel, error) {
	brc, err := repo.getBoostRequestChannel("WHERE guild_id = ? AND frontend_channel_id = ?", guildID, frontendChannelID)
	return brc, err
}

func (repo dbRepository) getBoostRequestChannel(where string, args ...interface{}) (*BoostRequestChannel, error) {
	row := repo.db.QueryRow(
		`SELECT id, guild_id, frontend_channel_id, backend_channel_id, uses_buyer_message, notifies_buyer
			FROM boost_request_channel `+where,
		args...,
	)
	brc := BoostRequestChannel{}
	var usesBuyerMessage, notifiesBuyer int
	err := row.Scan(&brc.ID, &brc.GuildID, &brc.FrontendChannelID, &brc.BackendChannelID, &usesBuyerMessage, &notifiesBuyer)
	if err == sql.ErrNoRows {
		return nil, ErrBoostRequestChannelNotFound
	}
	if err != nil {
		return nil, err
	}
	brc.UsesBuyerMessage = usesBuyerMessage != 0
	brc.NotifiesBuyer = notifiesBuyer != 0

	return &brc, nil
}
