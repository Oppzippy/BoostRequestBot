package repository

import "time"

type AutoSignupSession struct {
	GuildID      string
	AdvertiserID string
	ExpiresAt    time.Time
}
