package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func (repo *dbRepository) GetBoostRequestByBackendMessageID(backendChannelID, backendMessageID string) (*repository.BoostRequest, error) {
	return repo.getBoostRequest(
		`WHERE br.deleted_at IS NULL AND br.id = (
			SELECT
				brbm.boost_request_id
			FROM
				boost_request_backend_message AS brbm
			WHERE
				brbm.channel_id = ? AND
				brbm.message_id = ?
		)`,
		backendChannelID,
		backendMessageID,
	)
}

func (repo *dbRepository) GetUnresolvedBoostRequests() ([]*repository.BoostRequest, error) {
	// TODO don't load really old boost requests
	boostRequests, err := repo.getBoostRequests("WHERE br.resolved_at IS NULL AND br.deleted_at IS NULL")
	return boostRequests, err
}

func (repo *dbRepository) GetBoostRequestById(guildID string, boostRequestID uuid.UUID) (*repository.BoostRequest, error) {
	br, err := repo.getBoostRequest(`
		WHERE
			br.guild_id = ?
			AND br.external_id = ?
			AND br.deleted_at IS NULL`,
		guildID,
		boostRequestID,
	)
	return br, err
}

func (repo *dbRepository) getBoostRequest(where string, args ...interface{}) (*repository.BoostRequest, error) {
	boostRequests, err := repo.getBoostRequests(where, args...)
	if err != nil {
		return nil, fmt.Errorf("extracting single request: %w", err)
	}

	if len(boostRequests) == 0 {
		return nil, repository.ErrNoResults
	}
	return boostRequests[0], nil
}

func (repo *dbRepository) getBoostRequests(where string, args ...interface{}) ([]*repository.BoostRequest, error) {
	row, err := repo.db.Query(`
		SELECT
			br.id, br.external_id, br.guild_id, br.backend_channel_id, br.requester_id, br.advertiser_id, br.message,
			br.embed_fields, br.price, br.created_at, br.resolved_at, br.name_visibility, br.collect_users_only,
			brc.id, brc.guild_id, brc.frontend_channel_id, brc.backend_channel_id, brc.uses_buyer_message, brc.skips_buyer_dm
		FROM
			boost_request br
		LEFT JOIN boost_request_channel brc ON
			br.boost_request_channel_id = brc.id
		`+where,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	boostRequests := make([]*repository.BoostRequest, 0, 1) // Optimize for the common case of a specific boost request being selected

	for row.Next() {
		br, err := repo.unmarshalBoostRequest(row)
		if err != nil {
			return nil, err
		}
		// TODO do this in 2 queries rather than n+1
		preferredAdvertiserIDs, err := repo.getPreferredAdvertisers(br)
		if err != nil {
			return nil, err
		}
		br.PreferredAdvertiserIDs = preferredAdvertiserIDs

		// TODO another n+1
		backendMessages, err := repo.getBoostRequestBackendMessages(br.ID)
		if err != nil {
			return nil, err
		}
		br.BackendMessages = backendMessages

		boostRequests = append(boostRequests, br)
	}

	return boostRequests, nil
}

type scannable interface {
	Scan(dest ...interface{}) error
}

func (repo *dbRepository) unmarshalBoostRequest(row scannable) (*repository.BoostRequest, error) {
	var (
		br               repository.BoostRequest
		advertiserID     sql.NullString
		resolvedAt       sql.NullTime
		embedFieldsJSON  sql.NullString
		price            sql.NullInt64
		nameVisibility   string
		collectUsersOnly int

		brcID                sql.NullInt64
		brcGuildID           sql.NullString
		brcFrontendChannelID sql.NullString
		brcBackendChannelID  sql.NullString
		brcUsesBuyerMessage  sql.NullBool
		brcSkipsBuyerDM      sql.NullBool
	)
	err := row.Scan(
		&br.ID, &br.ExternalID, &br.GuildID, &br.BackendChannelID, &br.RequesterID, &advertiserID, &br.Message,
		&embedFieldsJSON, &price, &br.CreatedAt, &resolvedAt, &nameVisibility, &collectUsersOnly,
		&brcID, &brcGuildID, &brcFrontendChannelID, &brcBackendChannelID, &brcUsesBuyerMessage, &brcSkipsBuyerDM,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNoResults
		}
		return nil, err
	}
	if embedFieldsJSON.Valid {
		err := json.Unmarshal([]byte(embedFieldsJSON.String), &br.EmbedFields)
		if err != nil {
			log.Printf("Error parsing embed field json: %v", err)
		}
	}
	br.AdvertiserID = advertiserID.String
	br.IsResolved = resolvedAt.Valid
	br.Price = price.Int64
	br.NameVisibility = repository.NameVisibilitySettingFromString(nameVisibility)
	br.CollectUsersOnly = collectUsersOnly != 0
	if brcID.Valid {
		br.Channel = &repository.BoostRequestChannel{
			ID:                brcID.Int64,
			GuildID:           brcGuildID.String,
			FrontendChannelID: brcFrontendChannelID.String,
			BackendChannelID:  brcBackendChannelID.String,
			UsesBuyerMessage:  brcUsesBuyerMessage.Bool,
			SkipsBuyerDM:      brcUsesBuyerMessage.Bool,
		}
	}
	return &br, nil
}

// Inserts the boost request into the database and updates the ID field to match the newly inserted row's id
func (repo *dbRepository) InsertBoostRequest(br *repository.BoostRequest) error {
	embedFieldsJSON := sql.NullString{}
	if br.EmbedFields != nil {
		embedFieldsJSONBytes, err := json.Marshal(br.EmbedFields)
		if err != nil {
			log.Printf("Error marshalling embed field json: %v", err)
		} else {
			embedFieldsJSON.String = string(embedFieldsJSONBytes)
			embedFieldsJSON.Valid = true
		}
	}
	channelID := sql.NullInt64{}
	if br.Channel != nil {
		channelID.Int64 = br.Channel.ID
		channelID.Valid = true
	}
	var collectUsersOnly int
	if br.CollectUsersOnly {
		collectUsersOnly = 1
	}

	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Exec(
		`INSERT INTO boost_request
			(
				external_id, boost_request_channel_id, guild_id, backend_channel_id, requester_id, advertiser_id, message, embed_fields,
				price, created_at, name_visibility, collect_users_only
			)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		br.ExternalID,
		channelID,
		br.GuildID,
		br.BackendChannelID,
		br.RequesterID,
		sql.NullString{
			String: br.AdvertiserID,
			Valid:  br.AdvertiserID != "",
		},
		br.Message,
		embedFieldsJSON,
		sql.NullInt64{
			Int64: br.Price,
			Valid: br.Price != 0,
		},
		br.CreatedAt,
		br.NameVisibility.String(),
		collectUsersOnly,
	)
	err = rollbackIfErr(tx, err)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	err = rollbackIfErr(tx, err)
	if err != nil {
		return err
	}

	err = rollbackIfErr(tx, repo.updatePreferredAdvertisers(tx, id, br.PreferredAdvertiserIDs))
	if err != nil {
		return err
	}

	err = rollbackIfErr(tx, repo.updateBoostRequestBackendMessages(tx, id, br.BackendMessages))
	if err != nil {
		return err
	}

	br.ID = id
	err = tx.Commit()
	return err
}

func (repo *dbRepository) ResolveBoostRequest(br *repository.BoostRequest) error {
	var resolvedAt *time.Time
	if br.IsResolved {
		resolvedAt = &br.ResolvedAt
	}

	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`UPDATE boost_request SET
			advertiser_id = ?,
			resolved_at = ?
			WHERE id = ?`,
		br.AdvertiserID,
		resolvedAt,
		br.ID,
	)
	err = rollbackIfErr(tx, err)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func (repo *dbRepository) DeleteBoostRequest(br *repository.BoostRequest) error {
	_, err := repo.db.Exec(`
		UPDATE
			boost_request
		SET
			deleted_at = ?
		WHERE
			id = ?
		AND
			resolved_at IS NULL AND
			deleted_at IS NULL`,
		time.Now(),
		br.ID,
	)
	return err
}

func (repo *dbRepository) InsertBoostRequestDelayedMessage(br *repository.BoostRequest, delayedMessage *repository.DelayedMessage) error {
	_, err := repo.db.Exec(`
		INSERT INTO
			boost_request_delayed_message (
				boost_request_id,
				delayed_message_id
			)
		VALUES (?, ?)`,
		br.ID,
		delayedMessage.ID,
	)
	return err
}

func (repo *dbRepository) GetBoostRequestDelayedMessageIDs(br *repository.BoostRequest) ([]int64, error) {
	rows, err := repo.db.Query(`
		SELECT
			delayed_message_id
		FROM
			boost_request_delayed_message
		WHERE
			boost_request_id = ?`,
		br.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]int64, 0, 5)
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
