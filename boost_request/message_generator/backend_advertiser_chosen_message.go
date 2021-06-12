package message_generator

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type BackendAdvertiserChosenMessage struct {
	localizer           *i18n.Localizer
	boostRequest        *repository.BoostRequest
	advertiserAvatarURL string
}

func NewBackendAdvertiserChosenMessage(
	localizer *i18n.Localizer, br *repository.BoostRequest, advertiserAvatarURL string,
) *BackendAdvertiserChosenMessage {
	return &BackendAdvertiserChosenMessage{
		localizer:           localizer,
		boostRequest:        br,
		advertiserAvatarURL: advertiserAvatarURL,
	}
}

func (m *BackendAdvertiserChosenMessage) Message() (*discordgo.MessageSend, error) {
	return &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Color: 0xFF0000,
			Title: m.localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "AdvertiserSelected",
					Other: "An advertiser has been selected.",
				},
			}),
			Description: m.localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "AdvertiserWillHandleBoostRequest",
					Other: "{{.Advertiser}} will handle the following boost request.",
				},
				TemplateData: map[string]string{
					"Advertiser": fmt.Sprintf("<@%s>", m.boostRequest.AdvertiserID),
				},
			}),
			Fields: formatBoostRequest(m.localizer, m.boostRequest),
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: m.advertiserAvatarURL,
			},
		},
	}, nil
}

func formatBoostRequest(localizer *i18n.Localizer, br *repository.BoostRequest) []*discordgo.MessageEmbedField {
	var fields []*discordgo.MessageEmbedField
	if br.EmbedFields != nil {
		fields = repository.ToDiscordEmbedFields(br.EmbedFields)
	} else {
		fields = []*discordgo.MessageEmbedField{
			{
				Name: localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:  "BoostRequest",
						One: "Boost Request",
					},
					PluralCount: 1,
				}),
				Value: br.Message,
			},
		}
	}
	return fields
}
