package boost_request

import (
	"database/sql"
	"errors"
)

var ErrBoostRequestNotFound = errors.New("boost request not found")

func GetBoostRequestByBackendMessageID(db *sql.DB, backendChannelID, backendMessageID string) (*BoostRequest, error) {
	return getBoostRequest(
		db,
		"WHERE brc.backend_channel_id = ? AND br.backend_message_id = ?",
		backendChannelID,
		backendMessageID,
	)
}

func getBoostRequest(db *sql.DB, where string, args ...interface{}) (*BoostRequest, error) {
	row := db.QueryRow(`SELECT
		br.id, br.requester_id, br.advertiser_id, br.backend_message_id, br.message,
		brc.id, brc.guild_id, brc.frontend_channel_id, brc.backend_channel_id, brc.uses_buyer_message, brc.notifies_buyer
		FROM boost_request br
		INNER JOIN boost_request_channel brc ON br.boost_request_channel_id = brc.id `+where,
		args...,
	)
	if row.Err() == sql.ErrNoRows {
		return nil, ErrBoostRequestNotFound
	}

	var br BoostRequest
	var brc BoostRequestChannel
	row.Scan(
		&br.ID, &br.RequesterID, &br.RequesterID, &br.AdvertiserID, &br.BackendMessageID, &br.Message,
		&brc.ID, &brc.GuildID, &brc.FrontendChannelID, &brc.BackendChannelID, &brc.UsesBuyerMessage, &brc.NotifiesBuyer,
	)
	br.Channel = &brc
	return &br, nil
}
