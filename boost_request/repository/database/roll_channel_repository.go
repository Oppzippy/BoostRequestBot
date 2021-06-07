package database

import (
	"database/sql"

	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func (repo *dbRepository) GetRollChannel(guildID string) (channelID string, err error) {
	row := repo.db.QueryRow(
		"SELECT channel_id FROM roll_channel WHERE guild_id = ?",
		guildID,
	)
	err = row.Scan(&channelID)
	if err == sql.ErrNoRows {
		return "", repository.ErrNoResults
	}
	if err != nil {
		return "", err
	}
	return channelID, nil
}

func (repo *dbRepository) InsertRollChannel(guildID string, channelID string) error {
	_, err := repo.db.Exec(
		`INSERT INTO roll_channel (guild_id, channel_id) 
			VALUES (?, ?)
			ON DUPLICATE KEY UPDATE
				channel_id = VALUES(channel_id)`,
		guildID,
		channelID,
	)
	return err
}

func (repo *dbRepository) DeleteRollChannel(guildID string) error {
	_, err := repo.db.Exec("DELETE FROM roll_channel WHERE guild_id = ?", guildID)
	return err
}
