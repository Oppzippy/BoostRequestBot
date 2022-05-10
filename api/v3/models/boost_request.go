package models

import (
	"time"

	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type BoostRequest struct {
	ID                    string                `json:"id"`
	RequesterID           string                `json:"requesterId"`
	BackendChannelID      string                `json:"backendChannelId"`
	IsClaimed             bool                  `json:"isClaimed"`
	AdvertiserID          string                `json:"advertiserId,omitempty"`
	Message               string                `json:"message"`
	Price                 int64                 `json:"price,string,omitempty"`
	PreferredClaimerIDs   []string              `json:"preferredClaimerIds,omitempty"`
	AdditionalEmbedFields []*MessageEmbedField  `json:"additionalEmbedFields,omitempty"`
	CreatedAt             string                `json:"createdAt"`
	ClaimedAt             string                `json:"claimedAt,omitempty"`
	NameVisibility        NameVisibilitySetting `json:"nameVisibility"`
	DontPickClaimer       bool                  `json:"dontPickClaimer"`
}

type BoostRequestPartial struct {
	RequesterID           string                `json:"requesterId" validate:"required,numeric"`
	BackendChannelID      string                `json:"backendChannelId" validate:"required,numeric"`
	Message               string                `json:"message" validate:"required,max=1024"`
	Price                 int64                 `json:"price,string,omitempty" validate:"gte=0"`
	PreferredClaimerIDs   []string              `json:"preferredClaimerIds,omitempty" validate:"dive,required,numeric"`
	AdditionalEmbedFields []*MessageEmbedField  `json:"additionalEmbedFields,omitempty" validate:"dive"`
	NameVisibility        NameVisibilitySetting `json:"nameVisibility"`
	DontPickClaimer       bool                  `json:"dontPickClaimer"`
}

type MessageEmbedField struct {
	Name   string `json:"name" validate:"required,max=256"`
	Value  string `json:"value" validate:"required,max=1024"`
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

	nameVisibility := NameVisibilitySettingFromString(br.NameVisibility.String())

	return &BoostRequest{
		ID:                    br.ExternalID.String(),
		RequesterID:           br.RequesterID,
		BackendChannelID:      br.BackendChannelID,
		IsClaimed:             br.IsResolved,
		AdvertiserID:          br.AdvertiserID,
		Message:               br.Message,
		Price:                 br.Price,
		PreferredClaimerIDs:   preferredClaimerIds,
		AdditionalEmbedFields: embedFields,
		CreatedAt:             br.CreatedAt.Format(time.RFC3339),
		ClaimedAt:             claimedAt,
		NameVisibility:        nameVisibility,
		DontPickClaimer:       br.CollectUsersOnly,
	}
}
