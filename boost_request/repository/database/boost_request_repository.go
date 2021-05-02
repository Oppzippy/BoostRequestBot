package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func (repo *dbRepository) GetBoostRequestByBackendMessageID(backendChannelID, backendMessageID string) (*repository.BoostRequest, error) {
	return repo.getBoostRequest(
		"WHERE brc.backend_channel_id = ? AND br.backend_message_id = ? AND br.deleted_at IS NULL",
		backendChannelID,
		backendMessageID,
	)
}

func (repo *dbRepository) GetUnresolvedBoostRequests() ([]*repository.BoostRequest, error) {
	// TODO don't load really old boost requests
	boostRequests, err := repo.getBoostRequests("WHERE resolved_at IS NULL")
	return boostRequests, err
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
	row, err := repo.db.Query(
		`SELECT
			br.id, br.requester_id, br.advertiser_id, br.backend_message_id, br.message, br.embed_fields, br.created_at, br.resolved_at,
			brc.id, brc.guild_id, brc.frontend_channel_id, brc.backend_channel_id, brc.uses_buyer_message, brc.skips_buyer_dm
		FROM boost_request br
			INNER JOIN boost_request_channel brc ON br.boost_request_channel_id = brc.id
		`+where,
		args...,
	)
	// TODO fetch role discounts

	if err != nil {
		return nil, err
	}

	boostRequests := make([]*repository.BoostRequest, 0, 1) // Optimize for the common case of a specific boost request being selected

	for row.Next() {
		br, err := repo.unmarshalBoostRequest(row)
		if err != nil {
			return nil, err
		}
		boostRequests = append(boostRequests, br)
	}

	return boostRequests, nil
}

type scannable interface {
	Scan(dest ...interface{}) error
}

func (repo *dbRepository) unmarshalBoostRequest(row scannable) (*repository.BoostRequest, error) {
	var br repository.BoostRequest
	var brc repository.BoostRequestChannel
	var advertiserID sql.NullString
	var resolvedAt sql.NullTime
	var embedFieldsJSON sql.NullString
	err := row.Scan(
		&br.ID, &br.RequesterID, &advertiserID, &br.BackendMessageID, &br.Message, &embedFieldsJSON, &br.CreatedAt, &resolvedAt,
		&brc.ID, &brc.GuildID, &brc.FrontendChannelID, &brc.BackendChannelID, &brc.UsesBuyerMessage, &brc.SkipsBuyerDM,
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
	br.Channel = brc
	br.AdvertiserID = advertiserID.String
	br.IsResolved = resolvedAt.Valid
	return &br, nil
}

// Inserts the boost request into the database and updates the ID field to match the newly inserted row's id
func (repo *dbRepository) InsertBoostRequest(br *repository.BoostRequest) error {
	var advertiserID *string
	if br.AdvertiserID != "" {
		advertiserID = &br.AdvertiserID
	}
	var embedFieldsJSON *string
	if br.EmbedFields != nil {
		embedFieldsJSONBytes, err := json.Marshal(br.EmbedFields)
		if err != nil {
			log.Printf("Error marshalling embed field json: %v", err)
		} else {
			s := string(embedFieldsJSONBytes)
			embedFieldsJSON = &s
		}
	}
	res, err := repo.db.Exec(
		`INSERT INTO boost_request
			(boost_request_channel_id, requester_id, advertiser_id, backend_message_id, message, embed_fields, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
		br.Channel.ID,
		br.RequesterID,
		advertiserID,
		br.BackendMessageID,
		br.Message,
		embedFieldsJSON,
		br.CreatedAt,
	)
	// TODO insert role discounts if they're set
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
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return rbErr
		}
		return err
	}

	err = insertRoleDiscounts(br, tx)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return rbErr
		}
		return err
	}

	err = tx.Commit()
	return err
}

func insertRoleDiscounts(br *repository.BoostRequest, tx *sql.Tx) error {
	numRoleDiscounts := len(br.RoleDiscounts)
	args := make([]interface{}, 0, numRoleDiscounts+1)
	args = append(args, br.ID)
	if br.RoleDiscounts != nil {
		for _, rd := range br.RoleDiscounts {
			args = append(args, rd.ID)
		}
	}
	_, err := tx.Exec(
		`DELETE FROM boost_request_role_discount
		WHERE
			boost_request_id = ?
			AND role_discount_id NOT IN `+SQLSet(numRoleDiscounts),
		args...,
	)
	if err != nil {
		return err
	}
	if numRoleDiscounts == 0 {
		return nil
	}

	args = make([]interface{}, 0, numRoleDiscounts*2)
	for _, rd := range br.RoleDiscounts {
		args = append(args, br.ID, rd.ID)
	}

	_, err = tx.Exec(
		`INSERT IGNORE INTO boost_request_role_discount (
			boost_request_id, role_discount_id
		) VALUES `+SQLSets(2, numRoleDiscounts),
		args...,
	)
	return err
}
