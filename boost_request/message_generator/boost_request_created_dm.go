package message_generator

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type BoostRequestCreatedDM struct {
	localizer    *i18n.Localizer
	boostRequest *repository.BoostRequest
	userProvider userProvider
}

func NewBoostRequestCreatedDM(
	localizer *i18n.Localizer, userProvider userProvider, br *repository.BoostRequest,
) *BoostRequestCreatedDM {
	return &BoostRequestCreatedDM{
		localizer:    localizer,
		boostRequest: br,
		userProvider: userProvider,
	}
}

func (m *BoostRequestCreatedDM) Message() (*discordgo.MessageSend, error) {
	requester, err := m.userProvider.User(m.boostRequest.RequesterID)
	if err != nil {
		return nil, err
	}
	return &discordgo.MessageSend{
		Content: "Please wait while we find an advertiser to complete your request.",
		Embed: &discordgo.MessageEmbed{
			Title: "Huokan Boosting Community Boost Request",
			Author: &discordgo.MessageEmbedAuthor{
				Name: requester.String(),
			},
			Description: m.boostRequest.Message,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: requester.AvatarURL(""),
			},
		},
	}, nil
}
