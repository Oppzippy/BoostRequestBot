package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func (repo *dbRepository) GetAPIKey(key string) (*repository.APIKey, error) {
	row := repo.db.QueryRow("SELECT id, `key`, guild_id FROM api_key WHERE `key` = ?", key)

	var apiKey repository.APIKey
	err := row.Scan(&apiKey.ID, &apiKey.Key, &apiKey.GuildID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNoResults
		}
		return nil, err
	}
	return &apiKey, nil
}

func (repo *dbRepository) NewAPIKey(guildID string) (*repository.APIKey, error) {
	key, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	row, err := repo.db.Exec(
		"INSERT INTO api_key (`key`, guild_id, created_at) VALUES (?, ?, ?)",
		key.String(),
		guildID,
		time.Now().UTC(),
	)
	if err != nil {
		return nil, err
	}
	id, err := row.LastInsertId()
	if err != nil {
		return nil, err
	}
	apiKey := repository.APIKey{
		ID:      id,
		Key:     key.String(),
		GuildID: guildID,
	}
	return &apiKey, nil
}
