package messages

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages/partials"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type BoostRequestCreatedDM struct {
	localizer    *i18n.Localizer
	boostRequest *repository.BoostRequest
	userProvider userProvider
	embedPartial *partials.BoostRequestEmbedTemplate
}

func NewBoostRequestCreatedDM(
	localizer *i18n.Localizer, userProvider userProvider, br *repository.BoostRequest,
) *BoostRequestCreatedDM {
	return &BoostRequestCreatedDM{
		localizer:    localizer,
		boostRequest: br,
		userProvider: userProvider,
		embedPartial: partials.NewBoostRequestEmbedTemplate(localizer, br),
	}
}

func (m *BoostRequestCreatedDM) Message() (*discordgo.MessageSend, error) {
	requester, err := m.userProvider.User(m.boostRequest.RequesterID)
	if err != nil {
		return nil, err
	}
	var content string
	if len(m.boostRequest.PreferredAdvertiserIDs) == 0 {
		content = m.localizer.MustLocalize(
			&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "BoostRequestPleaseWait",
					Other: "Please wait while we find a booster to complete your request.",
				},
			},
		)
	} else {
		content = m.localizer.MustLocalize(
			&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "BoostRequestPleaseWaitPreferredClaimer",
					Other: "Please wait for your preferred claimer to claim the boost request. If you wish to remove your preference and accept any claimer, you may use the button below. ",
				},
			},
		)
	}

	embed, err := m.embedPartial.Embed(partials.BoostRequestEmbedConfiguration{
		PreferredAdvertisers: true,
		Description:          content,
		Price:                true,
	})
	if err != nil {
		return nil, err
	}
	embed.Title = m.localizer.MustLocalize(
		&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "HuokanBoostRequest",
				Other: "Huokan Boosting Community Boost Request",
			},
		},
	)
	embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: requester.AvatarURL(""),
	}
	embed.Author = &discordgo.MessageEmbedAuthor{
		Name: requester.String(),
	}

	var components []discordgo.MessageComponent
	if len(m.boostRequest.PreferredAdvertiserIDs) != 0 {
		components = m.components()
	}

	return &discordgo.MessageSend{
		Embed:      embed,
		Components: components,
	}, nil
}

func (m *BoostRequestCreatedDM) components() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label: m.localizer.MustLocalize(&i18n.LocalizeConfig{
						DefaultMessage: &i18n.Message{
							ID:    "RemoveClaimerPreference",
							Other: "Remove Claimer Preference",
						},
					}),
					Style: discordgo.PrimaryButton,
					CustomID: fmt.Sprintf(
						"removeAdvertiserPreference:%s:%s",
						m.boostRequest.GuildID,
						m.boostRequest.ExternalID.String(),
					),
				},
			},
		},
	}
}
