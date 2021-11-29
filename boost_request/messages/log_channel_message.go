package messages

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type LogChannelMessage struct {
	localizer    *i18n.Localizer
	boostRequest *repository.BoostRequest
	userProvider userProvider
}

func NewLogChannelMessage(
	localizer *i18n.Localizer, up userProvider, br *repository.BoostRequest,
) *LogChannelMessage {
	return &LogChannelMessage{
		localizer:    localizer,
		boostRequest: br,
		userProvider: up,
	}
}

func (m *LogChannelMessage) Message() (*discordgo.MessageSend, error) {
	if m.boostRequest.Channel != nil && m.boostRequest.Channel.UsesBuyerMessage {
		return nil, errors.New("UsesBuyerMessage logging not implemented")
	}
	user, err := m.userProvider.User(m.boostRequest.RequesterID)
	if err != nil {
		return nil, err
	}

	return &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Color: 0x0000FF,
			Title: m.localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "NewBoostRequest",
					One:   "New Boost Request",
					Other: "New Boost Requests",
				},
				PluralCount: 1,
			}),
			Description: m.boostRequest.Message,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name: m.localizer.MustLocalize(&i18n.LocalizeConfig{
						DefaultMessage: &i18n.Message{
							ID:    "RequestedBy",
							Other: "Requested By",
						},
					}),
					Value: user.Mention() + " " + user.String(),
				},
			},
		},
	}, nil
}
