package repository

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

type MessageEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type BoostRequest struct {
	ID int64
	// TODO generate uuids for all nulls and make the column not null
	ExternalID             *uuid.UUID
	Channel                *BoostRequestChannel
	GuildID                string
	BackendChannelID       string
	RequesterID            string
	AdvertiserID           string
	BackendMessages        []*BoostRequestBackendMessage
	Message                string
	EmbedFields            []*MessageEmbedField
	Price                  int64
	PreferredAdvertiserIDs map[string]struct{}
	CreatedAt              time.Time
	IsResolved             bool
	ResolvedAt             time.Time
	NameVisibility         NameVisibilitySetting
	CollectUsersOnly       bool
}

func FromDiscordEmbedFields(fields []*discordgo.MessageEmbedField) []*MessageEmbedField {
	if fields == nil {
		return nil
	}
	newFields := make([]*MessageEmbedField, len(fields))
	for i, field := range fields {
		newFields[i] = &MessageEmbedField{
			Name:   field.Name,
			Value:  field.Value,
			Inline: field.Inline,
		}
	}
	return newFields
}

func ToDiscordEmbedFields(fields []*MessageEmbedField) []*discordgo.MessageEmbedField {
	if fields == nil {
		return nil
	}
	newFields := make([]*discordgo.MessageEmbedField, len(fields))
	for i, field := range fields {
		newFields[i] = &discordgo.MessageEmbedField{
			Name:   field.Name,
			Value:  field.Value,
			Inline: field.Inline,
		}
	}
	return newFields
}
