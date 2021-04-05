package repository

import (
	"database/sql"
	"errors"
	"time"
)

type BoostRequestRepository interface {
	GetBoostRequestByBackendMessageID(backendChannelID, backendMessageID string) (*BoostRequest, error)
	InsertBoostRequest(br *BoostRequest) error
	ResolveBoostRequest(br *BoostRequest) error
}

var ErrBoostRequestNotFound = errors.New("boost request not found")

func (repo dbRepository) GetBoostRequestByBackendMessageID(backendChannelID, backendMessageID string) (*BoostRequest, error) {
	return repo.getBoostRequest(
		"WHERE brc.backend_channel_id = ? AND br.backend_message_id = ? AND deleted_at IS NULL",
		backendChannelID,
		backendMessageID,
	)
}

func (repo dbRepository) getBoostRequest(where string, args ...interface{}) (*BoostRequest, error) {
	row := repo.db.QueryRow(`SELECT
		br.id, br.requester_id, br.advertiser_id, br.backend_message_id, br.message,
		brc.id, brc.guild_id, brc.frontend_channel_id, brc.backend_channel_id, brc.uses_buyer_message, brc.skips_buyer_dm
		FROM boost_request br
		INNER JOIN boost_request_channel brc ON br.boost_request_channel_id = brc.id `+where,
		args...,
	)

	var br BoostRequest
	var brc BoostRequestChannel
	br.Channel = &brc
	err := row.Scan(
		&br.ID, &br.RequesterID, &br.RequesterID, &br.AdvertiserID, &br.BackendMessageID, &br.Message,
		&brc.ID, &brc.GuildID, &brc.FrontendChannelID, &brc.BackendChannelID, &brc.UsesBuyerMessage, &brc.SkipsBuyerDM,
	)
	if err == sql.ErrNoRows {
		return nil, ErrBoostRequestNotFound
	} else if err != nil {
		return nil, err
	}
	return &br, nil
}

// Inserts the boost request into the database and updates the ID field to match the newly inserted row's id
func (repo dbRepository) InsertBoostRequest(br *BoostRequest) error {
	var advertiserID *string = nil
	if br.AdvertiserID != "" {
		advertiserID = &br.AdvertiserID
	}
	res, err := repo.db.Exec(
		`INSERT INTO boost_request
			(boost_request_channel_id, requester_id, backend_message_id, message, created_at)
			VALUES (?, ?, ?, ?, ?)`,
		br.Channel.ID,
		br.RequesterID,
		advertiserID,
		br.BackendMessageID,
		br.Message,
		br.CreatedAt,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	br.ID = id
	return nil
}

func (repo dbRepository) ResolveBoostRequest(br *BoostRequest) error {
	var resolvedAt *string = nil
	if br.IsResolved {
		timeStr := br.ResolvedAt.UTC().Format(time.RFC3339)
		resolvedAt = &timeStr
	}
	_, err := repo.db.Exec(
		`UPDATE boost_request SET
			advertiser_id = ?
			resolved_at = ?
			WHERE id = ?`,
		br.AdvertiserID,
		resolvedAt,
		br.ID,
	)
	return err
}
