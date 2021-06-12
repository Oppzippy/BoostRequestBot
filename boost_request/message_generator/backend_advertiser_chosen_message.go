package message_generator

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type BackendAdvertiserChosenMessage struct {
	localizer    *i18n.Localizer
	boostRequest *repository.BoostRequest
	userProvider userProvider
}

func NewBackendAdvertiserChosenMessage(
	localizer *i18n.Localizer, up userProvider, br *repository.BoostRequest,
) *BackendAdvertiserChosenMessage {
	return &BackendAdvertiserChosenMessage{
		localizer:    localizer,
		boostRequest: br,
		userProvider: up,
	}
}

func (m *BackendAdvertiserChosenMessage) Message() (*discordgo.MessageSend, error) {
	advertiser, err := m.userProvider.User(m.boostRequest.AdvertiserID)
	if err != nil {
		return nil, err
	}

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
					Other: "{{.AdvertiserMention}} will handle the following boost request.",
				},
				TemplateData: map[string]string{
					"AdvertiserMention": fmt.Sprintf("<@%s>", m.boostRequest.AdvertiserID),
				},
			}),
			Fields: formatBoostRequest(m.localizer, m.boostRequest),
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: advertiser.AvatarURL(""),
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
