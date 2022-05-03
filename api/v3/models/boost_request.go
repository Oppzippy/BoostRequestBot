package models

import (
	"time"

	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type BoostRequest struct {
	ID                     string               `json:"id"`
	RequesterID            string               `json:"requesterId"`
	BackendChannelID       string               `json:"backendChannelId"`
	IsClaimed              bool                 `json:"isClaimed"`
	AdvertiserID           string               `json:"advertiserId,omitempty"`
	Message                string               `json:"message"`
	Price                  int64                `json:"price,string,omitempty"`
	PreferredAdvertiserIDs []string             `json:"preferredAdvertiserIds,omitempty"`
	AdditionalEmbedFields  []*MessageEmbedField `json:"additionalEmbedFields,omitempty"`
	CreatedAt              string               `json:"createdAt"`
	ClaimedAt              string               `json:"claimedAt,omitempty"`
}

type BoostRequestPartial struct {
	RequesterID           string               `json:"requesterId" validate:"required"`
	BackendChannelID      string               `json:"backendChannelId" validate:"required"`
	Message               string               `json:"message" validate:"required"`
	Price                 int64                `json:"price,string,omitempty"`
	PreferredClaimerIDs   []string             `json:"preferredClaimerIds,omitempty"`
	AdditionalEmbedFields []*MessageEmbedField `json:"additionalEmbedFields,omitempty" validate:"dive"`
}

type MessageEmbedField struct {
	Name   string `json:"name" validate:"required"`
	Value  string `json:"value" validate:"required"`
	Inline bool   `json:"inline"`
}

type BoostRequestBackendMessage struct {
	ChannelID string `json:"channelId" validate:"required"`
	MessageID string `json:"messageId" validate:"required"`
}

func FromRepositoryBoostRequest(br *repository.BoostRequest) *BoostRequest {
	preferredClaimerIds := make([]string, 0, len(br.PreferredAdvertiserIDs))
	if len(br.PreferredAdvertiserIDs) > 0 {
		for id := range br.PreferredAdvertiserIDs {
			preferredClaimerIds = append(preferredClaimerIds, id)
		}
	}

	embedFields := make([]*MessageEmbedField, len(br.EmbedFields))
	for i, embedField := range br.EmbedFields {
		embedFields[i] = &MessageEmbedField{
			Name:   embedField.Name,
			Value:  embedField.Value,
			Inline: embedField.Inline,
		}
	}

	var claimedAt string
	if !br.ResolvedAt.IsZero() {
		claimedAt = br.ResolvedAt.Format(time.RFC3339)
	}

	return &BoostRequest{
		ID:                     br.ExternalID.String(),
		RequesterID:            br.RequesterID,
		BackendChannelID:       br.BackendChannelID,
		IsClaimed:              br.IsResolved,
		AdvertiserID:           br.AdvertiserID,
		Message:                br.Message,
		Price:                  br.Price,
		PreferredAdvertiserIDs: preferredClaimerIds,
		AdditionalEmbedFields:  embedFields,
		CreatedAt:              br.CreatedAt.Format(time.RFC3339),
		ClaimedAt:              claimedAt,
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
