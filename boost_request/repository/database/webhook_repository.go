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
			guild_id = ?`,
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
			created_at = VALUES(created_at)`,
		webhook.GuildID,
		webhook.URL,
		time.Now().UTC(),
	)
	return err
}

func (repo *dbRepository) DeleteWebhook(webhook repository.Webhook) error {
	_, err := repo.db.Exec(`
		DELETE FROM
			webhook
		WHERE
			id = ?`,
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
			MAX(webhook_attempt.created_at)
		FROM
			webhook
		INNER JOIN webhook_queue ON
			webhook_queue.webhook_id = webhook.id
		LEFT JOIN webhook_attempt ON
			webhook_attempt.webhook_queue_id = webhook_queue.id
		WHERE
			webhook_attempt.id IS NULL OR
			webhook_attempt.status_code NOT BETWEEN 200 AND 299
		GROUP BY
			webhook.id,
			webhook_queue.id`,
	)
	if err != nil {
		return nil, err
	}
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
