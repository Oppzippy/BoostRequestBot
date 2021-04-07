package repository

import (
	"time"
)

type BoostRequest struct {
	ID               int64
	Channel          BoostRequestChannel
	RequesterID      string
	AdvertiserID     string
	BackendMessageID string
	Message          string
	CreatedAt        time.Time
	IsResolved       bool
	ResolvedAt       time.Time
}
