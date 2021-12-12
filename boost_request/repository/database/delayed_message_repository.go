package database

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func (repo *dbRepository) GetDelayedMessages() ([]*repository.DelayedMessage, error) {
	rows, err := repo.db.Query(`
		SELECT
			id,
			destination_id,
			destination_type,
			fallback_channel_id,
			message_json,
			send_at
		FROM
			delayed_message
		WHERE
			sent_at IS NULL AND
			deleted_at IS NULL AND
			send_at > ?`,
		// Don't bother with messages over an hour old
		time.Now().Add(-1*time.Hour),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]*repository.DelayedMessage, 0, 20)
	for rows.Next() {
		var (
			delayedMessage    repository.DelayedMessage
			messageJSON       string
			fallbackChannelID sql.NullString
		)
		err := rows.Scan(
			&delayedMessage.ID,
			&delayedMessage.DestinationID,
			&delayedMessage.DestinationType,
			&fallbackChannelID,
			&messageJSON,
			&delayedMessage.SendAt,
		)
		if err != nil {
			return nil, err
		}
		delayedMessage.FallbackChannelID = fallbackChannelID.String

		var message discordgo.MessageSend
		err = json.Unmarshal([]byte(messageJSON), &message)
		if err != nil {
			log.Printf("failed to unmarshal delayed message (id %v): %v", delayedMessage.ID, err)
			continue
		}
		delayedMessage.Message = &message

		messages = append(messages, &delayedMessage)
	}
	return messages, nil
}

func (repo *dbRepository) InsertDelayedMessage(delayedMessage *repository.DelayedMessage) error {
	messageJSON, err := json.Marshal(delayedMessage.Message)
	if err != nil {
		return err
	}
	res, err := repo.db.Exec(`
		INSERT INTO delayed_message (
			destination_id,
			destination_type,
			fallback_channel_id,
			message_json,
			send_at
		) VALUES (?, ?, ?, ?, ?)`,
		delayedMessage.DestinationID,
		delayedMessage.DestinationType,
		sql.NullString{
			String: delayedMessage.FallbackChannelID,
			Valid:  delayedMessage.FallbackChannelID != "",
		},
		messageJSON,
		delayedMessage.SendAt.UTC(),
	)
	if err != nil {
		return err
	}
	delayedMessage.ID, err = res.LastInsertId()
	return err
}

func (repo *dbRepository) DeleteDelayedMessage(id int64) error {
	_, err := repo.db.Exec(`
		UPDATE
			delayed_message
		SET
			deleted_at = ?
		WHERE
			id = ? AND
			deleted_at IS NULL`,
		time.Now().UTC(),
		id,
	)
	return err
}

func (repo *dbRepository) FlagDelayedMessageAsSent(message *repository.DelayedMessage) error {
	_, err := repo.db.Exec(`
		UPDATE
			delayed_message
		SET
			sent_at = ?
		WHERE
			id = ? AND
			sent_at IS NULL`,
		time.Now().UTC(),
		message.ID,
	)
	return err
}
