package repository

import "database/sql"

type LogChannelRepository interface {
	GetLogChannel(guildID string) (channelID string, err error)
	InsertLogChannel(guildID, channelID string) error
	DeleteLogChannel(guildID string) error
}

func (repo dbRepository) GetLogChannel(guildID string) (channelID string, err error) {
	row := repo.db.QueryRow("SELECT channel_id FROM log_channel WHERE guild_id = ?", guildID)
	err = row.Scan(&channelID)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	return channelID, nil
}

func (repo dbRepository) InsertLogChannel(guildID, channelID string) error {
	_, err := repo.db.Exec(
		`INSERT INTO log_channel (guild_id, channel_id)
			VALUES (?, ?)
			ON DUPLICATE KEY UPDATE
				channel_id = VALUES(channel_id)`,
		guildID,
		channelID,
	)
	return err
}

func (repo dbRepository) DeleteLogChannel(guildID string) error {
	_, err := repo.db.Exec("DELETE FROM log_channel WHERE guild_id = ?", guildID)
	return err
}
