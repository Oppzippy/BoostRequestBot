package models

import (
	"strconv"
	"time"

	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type BoostRequest struct {
	ID                     string            `json:"id"`
	RequesterID            string            `json:"requesterId"`
	IsAdvertiserSelected   bool              `json:"isAdvertiserSelected"`
	AdvertiserID           string            `json:"advertiserId,omitempty"`
	Message                string            `json:"message"`
	Price                  int64             `json:"price,string,omitempty"`
	Discount               int64             `json:"discount,string,omitempty"`
	AdvertiserCut          int64             `json:"advertiserCut,string,omitempty"`
	AdvertiserRoleCuts     map[string]string `json:"advertiserRoleCuts,omitempty"`
	PreferredAdvertiserIDs []string          `json:"preferredAdvertiserIds,omitempty"`
	CreatedAt              string            `json:"createdAt"`
	AdvertiserSelectedAt   string            `json:"advertiserSelectedAt,omitempty"`
}

type BoostRequestPartial struct {
	RequesterID            string            `json:"requesterId" validate:"required"`
	BackendChannelID       string            `json:"backendChannelId" validate:"required"`
	Message                string            `json:"message" validate:"required"`
	Price                  int64             `json:"price,string,omitempty"`
	AdvertiserCut          int64             `json:"advertiserCut,string,omitempty"`
	AdvertiserRoleCuts     map[string]string `json:"advertiserRoleCuts,omitempty"`
	Discount               int64             `json:"discount,string,omitempty"`
	PreferredAdvertiserIDs []string          `json:"preferredAdvertiserIds,omitempty"`
}

type BoostRequestBackendMessage struct {
	ChannelID string `json:"channelId" validate:"required"`
	MessageID string `json:"messageId" validate:"required"`
}

func FromRepositoryBoostRequest(br *repository.BoostRequest) *BoostRequest {
	roleCuts := make(map[string]string)
	if len(br.AdvertiserRoleCuts) > 0 {
		for roleID, cut := range br.AdvertiserRoleCuts {
			roleCuts[roleID] = strconv.FormatInt(cut, 10)
		}
	}

	preferredAdvertiserIDs := make([]string, 0, len(br.PreferredAdvertiserIDs))
	if len(br.PreferredAdvertiserIDs) > 0 {
		for id := range br.PreferredAdvertiserIDs {
			preferredAdvertiserIDs = append(preferredAdvertiserIDs, id)
		}
	}

	var advertiserSelectedAt string
	if !br.ResolvedAt.IsZero() {
		advertiserSelectedAt = br.ResolvedAt.Format(time.RFC3339)
	}

	return &BoostRequest{
		ID:                     br.ExternalID.String(),
		RequesterID:            br.RequesterID,
		IsAdvertiserSelected:   br.IsResolved,
		AdvertiserID:           br.AdvertiserID,
		Message:                br.Message,
		Price:                  br.Price,
		Discount:               br.Discount,
		AdvertiserCut:          br.AdvertiserCut,
		AdvertiserRoleCuts:     roleCuts,
		PreferredAdvertiserIDs: preferredAdvertiserIDs,
		CreatedAt:              br.CreatedAt.Format(time.RFC3339),
		AdvertiserSelectedAt:   advertiserSelectedAt,
	}
}

func FromRepositoryBackendMessages(messages []*repository.BoostRequestBackendMessage) []*BoostRequestBackendMessage {
	newMessages := make([]*BoostRequestBackendMessage, len(messages))
	for i, m := range messages {
		newMessages[i] = &BoostRequestBackendMessage{
			ChannelID: m.ChannelID,
			MessageID: m.MessageID,
		}
	}
	return newMessages
}
