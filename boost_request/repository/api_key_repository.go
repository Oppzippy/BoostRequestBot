package repository

import (
	"database/sql"
	"errors"
)

type ApiKeyRepository interface {
	GetAPIKey(key string) (*APIKey, error)
}

var ErrApiKeyNotFound = errors.New("api key not found")

func (repo *dbRepository) GetAPIKey(key string) (*APIKey, error) {
	row := repo.db.QueryRow("SELECT id, `key`, guild_id FROM api_key WHERE `key` = ?", key)

	var apiKey APIKey
	err := row.Scan(&apiKey.ID, &apiKey.Key, &apiKey.GuildID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrApiKeyNotFound
		}
		return nil, err
	}
	return &apiKey, nil
}
