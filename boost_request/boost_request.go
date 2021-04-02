package boost_request

import (
	"time"
)

type BoostRequest struct {
	ID               int
	Channel          *BoostRequestChannel
	RequesterID      string
	AdvertiserID     string
	BackendMessageID string
	Message          string
	CreatedAt        *time.Time
}
