package repository

import "time"

type Webhook struct {
	ID      int64
	GuildID string
	URL     string
}

type QueuedWebhookRequest struct {
	ID            int64
	Webhook       Webhook
	Body          string
	CreatedAt     time.Time
	LatestAttempt *time.Time
}

type WebhookAttempt struct {
	QueuedWebhookRequest QueuedWebhookRequest
	StatusCode           int
	CreatedAt            time.Time
}
