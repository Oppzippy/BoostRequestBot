package repository

import "time"

type AutoSignupSession struct {
	ID           int64
	GuildID      string
	AdvertiserID string
	ExpiresAt    time.Time
}
