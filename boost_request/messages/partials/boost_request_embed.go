package partials

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages/message_utils"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type BoostRequestEmbedTemplate struct {
	boostRequest *repository.BoostRequest
	localizer    *i18n.Localizer
}

type BoostRequestEmbedConfiguration struct {
	PreferredAdvertisers bool
	Description          string
	Price                bool
	ID                   bool
}

func NewBoostRequestEmbedTemplate(
	localizer *i18n.Localizer, br *repository.BoostRequest,
) *BoostRequestEmbedTemplate {
	return &BoostRequestEmbedTemplate{
		boostRequest: br,
		localizer:    localizer,
	}
}

func (m *BoostRequestEmbedTemplate) Embed(config BoostRequestEmbedConfiguration) (*discordgo.MessageEmbed, error) {
	embed := &discordgo.MessageEmbed{
		Title: m.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "BoostRequest",
				One:   "Boost Request",
				Other: "Boost Requests",
			},
			PluralCount: 1,
		}),
		Description: m.boostRequest.Message,
		Fields:      make([]*discordgo.MessageEmbedField, 0, 10),
		Color:       0x0000FF,
		Timestamp:   m.boostRequest.CreatedAt.UTC().Format(time.RFC3339),
	}

	if preferredAdvertisers := m.preferredAdvertisersField(); config.PreferredAdvertisers && preferredAdvertisers != nil {
		embed.Fields = append(embed.Fields, preferredAdvertisers)
	}
	if config.Description != "" {
		embed.Description = config.Description
		if message := m.messageField(); message != nil {
			embed.Fields = append(embed.Fields, message)
		}
	}
	if price := m.priceField(); config.Price && price != nil {
		embed.Fields = append(embed.Fields, price)
	}
	if len(m.boostRequest.EmbedFields) != 0 {
		embed.Fields = append(embed.Fields, repository.ToDiscordEmbedFields(m.boostRequest.EmbedFields)...)
	}
	if config.ID {
		embed.Footer = m.idFooter()
	}

	if len(embed.Fields) == 0 {
		embed.Fields = nil
	}

	return embed, nil
}

func (m *BoostRequestEmbedTemplate) preferredAdvertisersField() *discordgo.MessageEmbedField {
	if len(m.boostRequest.PreferredAdvertiserIDs) > 0 {
		mentions := make([]string, 0, len(m.boostRequest.PreferredAdvertiserIDs))
		for id := range m.boostRequest.PreferredAdvertiserIDs {
			mentions = append(mentions, fmt.Sprintf("<@%s>", id))
		}
		return &discordgo.MessageEmbedField{
			Name: m.localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "PreferredClaimer",
					One:   "Preferred Claimer",
					Other: "Preferred Claimers",
				},
				PluralCount: len(mentions),
			}),
			Value: strings.Join(mentions, " "),
		}
	}
	return nil
}

func (m *BoostRequestEmbedTemplate) messageField() *discordgo.MessageEmbedField {
	if m.boostRequest.Message != "" {
		return &discordgo.MessageEmbedField{
			Name: m.localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "Message",
					One:   "Message",
					Other: "Messages",
				},
				PluralCount: 1,
			}),
			Value: m.boostRequest.Message,
		}
	}
	return nil
}

func (m *BoostRequestEmbedTemplate) priceField() *discordgo.MessageEmbedField {
	if m.boostRequest.Price != 0 {
		return &discordgo.MessageEmbedField{
			Name:  "Price",
			Value: message_utils.FormatCopper(m.localizer, m.boostRequest.Price),
		}
	}
	return nil
}

func (m *BoostRequestEmbedTemplate) idFooter() *discordgo.MessageEmbedFooter {
	return &discordgo.MessageEmbedFooter{
		Text: m.boostRequest.ExternalID.String(),
	}
}
