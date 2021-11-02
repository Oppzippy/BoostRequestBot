package messages

import (
	"fmt"

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
	if len(m.boostRequest.PreferredAdvertiserIDs) > 0 {
		return m.preferredAdvertiserMessage()
	}
	return m.standardMessage()
}

func (m *BoostRequestCreatedDM) standardMessage() (*discordgo.MessageSend, error) {
	requester, err := m.userProvider.User(m.boostRequest.RequesterID)
	if err != nil {
		return nil, err
	}
	return &discordgo.MessageSend{
		Content: m.localizer.MustLocalize(
			&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "BoostRequestPleaseWait",
					Other: "Please wait while we find an advertiser to complete your request.",
				},
			},
		),
		Embed: &discordgo.MessageEmbed{
			Title: m.localizer.MustLocalize(
				&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "HuokanBoostRequest",
						Other: "Huokan Boosting Community Boost Request",
					},
				},
			),
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

func (m *BoostRequestCreatedDM) preferredAdvertiserMessage() (*discordgo.MessageSend, error) {
	requester, err := m.userProvider.User(m.boostRequest.RequesterID)
	if err != nil {
		return nil, err
	}
	return &discordgo.MessageSend{
		Content: m.localizer.MustLocalize(
			&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "BoostRequestPleaseWaitPreferredAdvertiser",
					Other: "Please wait for your preferred advertiser to claim the boost request. If you wish to remove your preference and accept any advertiser, you may use the button below. ",
				},
			},
		),
		Embed: &discordgo.MessageEmbed{
			Title: m.localizer.MustLocalize(
				&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "HuokanBoostRequest",
						Other: "Huokan Boosting Community Boost Request",
					},
				},
			),
			Author: &discordgo.MessageEmbedAuthor{
				Name: requester.String(),
			},
			Description: m.boostRequest.Message,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: requester.AvatarURL(""),
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label: m.localizer.MustLocalize(&i18n.LocalizeConfig{
							DefaultMessage: &i18n.Message{
								ID:    "RemoveAdvertiserPreference",
								Other: "Remove Advertiser Preference",
							},
						}),
						Style: discordgo.PrimaryButton,
						CustomID: fmt.Sprintf(
							"removeAdvertiserPreference:%s:%s",
							m.boostRequest.Channel.GuildID,
							m.boostRequest.ExternalID.String(),
						),
					},
				},
			},
		},
	}, nil
}
