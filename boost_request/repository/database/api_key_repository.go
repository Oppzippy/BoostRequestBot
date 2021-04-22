package database

import (
	"database/sql"

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
