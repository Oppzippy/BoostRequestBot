package database

import (
	"database/sql"
	"time"

	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func (repo *dbRepository) GetWebhook(guildId string) (repository.Webhook, error) {
	row := repo.db.QueryRow(`
		SELECT
			id, guild_id, webhook_url
		FROM
			webhook
		WHERE
			guild_id = ? AND
			deleted_at IS NULL`,
		guildId,
	)
	webhook := repository.Webhook{}
	err := row.Scan(&webhook.ID, &webhook.GuildID, &webhook.URL)
	if err == sql.ErrNoRows {
		return repository.Webhook{}, repository.ErrNoResults
	}
	if err != nil {
		return repository.Webhook{}, err
	}
	return webhook, nil
}

func (repo *dbRepository) InsertWebhook(webhook repository.Webhook) error {
	_, err := repo.db.Exec(`
		INSERT INTO webhook (
			guild_id, webhook_url, created_at
		) VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE
			webhook_url = VALUES(webhook_url),
			created_at = VALUES(created_at),
			deleted_at = NULL`,
		webhook.GuildID,
		webhook.URL,
		time.Now().UTC(),
	)
	return err
}

func (repo *dbRepository) DeleteWebhook(webhook repository.Webhook) error {
	_, err := repo.db.Exec(`
		UPDATE
			webhook
		SET
			deleted_at = ?
		WHERE
			id = ?`,
		time.Now(),
		webhook.ID,
	)
	return err
}

func (repo *dbRepository) InsertQueuedWebhook(webhook repository.Webhook, body string) error {
	_, err := repo.db.Exec(`
		INSERT INTO webhook_queue (
			webhook_id, body, created_at
		) VALUES (?, ?, ?)`,
		webhook.ID,
		body,
		time.Now(),
	)
	return err
}

func (repo *dbRepository) GetQueuedWebhooks() ([]*repository.QueuedWebhookRequest, error) {
	rows, err := repo.db.Query(`
		SELECT
			webhook.id,
			webhook.guild_id,
			webhook.webhook_url,
			webhook_queue.id,
			webhook_queue.body,
			webhook_queue.created_at,
			(SELECT webhook_attempt.created_at FROM webhook_attempt WHERE webhook_attempt.webhook_queue_id = webhook_queue.id LIMIT 1)	
		FROM
			webhook
		INNER JOIN webhook_queue ON
			webhook_queue.webhook_id = webhook.id
		WHERE
			webhook.deleted_at IS NULL AND
			webhook_queue.created_at < ? AND
			NOT EXISTS(
			    SELECT 1
			    FROM
			        webhook_attempt
			    WHERE
			        webhook_attempt.webhook_queue_id = webhook_queue.id AND
			        webhook_attempt.status_code BETWEEN 200 AND 299
			)`,
		time.Now().UTC().Add(time.Hour*24*7), // Give up after a week
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	queuedWebhooks := make([]*repository.QueuedWebhookRequest, 0)
	for rows.Next() {
		request := repository.QueuedWebhookRequest{}
		err := rows.Scan(
			&request.Webhook.ID,
			&request.Webhook.GuildID,
			&request.Webhook.URL,
			&request.ID,
			&request.Body,
			&request.CreatedAt,
			&request.LatestAttempt,
		)
		if err != nil {
			return nil, err
		}
		queuedWebhooks = append(queuedWebhooks, &request)
	}

	return queuedWebhooks, nil
}

func (repo *dbRepository) InsertWebhookAttempt(attempt repository.WebhookAttempt) error {
	_, err := repo.db.Exec(`
		INSERT INTO webhook_attempt (
			webhook_queue_id, status_code, created_at
		) VALUES (?, ?, ?)`,
		attempt.QueuedWebhookRequest.ID,
		attempt.StatusCode,
		attempt.CreatedAt,
	)
	return err
}
