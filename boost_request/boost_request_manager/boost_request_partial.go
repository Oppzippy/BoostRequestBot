package boost_request_manager

import (
	"github.com/oppzippy/BoostRequestBot/api/v3/models"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type BoostRequestPartial struct {
	GuildID                string
	BackendChannelID       string
	RequesterID            string
	Message                string
	EmbedFields            []*repository.MessageEmbedField
	PreferredAdvertiserIDs map[string]struct{}
	BackendMessageID       string
	Price                  int64
}

func FromModelBoostRequestPartial(br *models.BoostRequestPartial) (*BoostRequestPartial, error) {
	preferredAdvertiserIDs := make(map[string]struct{}, len(br.PreferredClaimerIDs))
	for _, advertiserID := range br.PreferredClaimerIDs {
		preferredAdvertiserIDs[advertiserID] = struct{}{}
	}

	embedFields := make([]*repository.MessageEmbedField, len(br.AdditionalEmbedFields))
	for i, embedField := range br.AdditionalEmbedFields {
		embedFields[i] = &repository.MessageEmbedField{
			Name:   embedField.Name,
			Value:  embedField.Value,
			Inline: embedField.Inline,
		}
	}

	brPartial := &BoostRequestPartial{
		RequesterID:            br.RequesterID,
		Message:                br.Message,
		PreferredAdvertiserIDs: preferredAdvertiserIDs,
		Price:                  br.Price,
		BackendChannelID:       br.BackendChannelID,
		EmbedFields:            embedFields,
	}

	return brPartial, nil
}
