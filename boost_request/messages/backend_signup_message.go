package messages

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages/partials"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type BackendSignupMessage struct {
	boostRequest      *repository.BoostRequest
	localizer         *i18n.Localizer
	discountFormatter *partials.DiscountFormatter
	embedPartial      *partials.BoostRequestEmbedPartial
}

func NewBackendSignupMessage(
	localizer *i18n.Localizer, df *partials.DiscountFormatter, br *repository.BoostRequest,
) *BackendSignupMessage {
	return &BackendSignupMessage{
		boostRequest:      br,
		localizer:         localizer,
		discountFormatter: df,
		embedPartial:      partials.NewBoostRequestEmbedPartial(localizer, df, br),
	}
}

func (m *BackendSignupMessage) Message() (*discordgo.MessageSend, error) {
	br := m.boostRequest
	embed, err := m.embedPartial.Embed(partials.BoostRequestEmbedConfiguration{
		Price:          true,
		AdvertiserCut:  true,
		Discount:       true,
		DiscountTotals: true,
	})
	if err != nil {
		return nil, err
	}
	embed.Title = m.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "NewBoostRequest",
			One:   "New Boost Request",
			Other: "New Boost Requests",
		},
		PluralCount: 1,
	})

	var preferredAdvertiserMentions string
	if len(br.PreferredAdvertiserIDs) > 0 {
		mentions := make([]string, len(br.PreferredAdvertiserIDs))
		for i, id := range br.PreferredAdvertiserIDs {
			mentions[i] = fmt.Sprintf("<@%s>", id)
		}
		preferredAdvertiserMentions = strings.Join(mentions, " ")
	}

	return &discordgo.MessageSend{
		Content: preferredAdvertiserMentions,
		Embed:   embed,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label: m.localizer.MustLocalize(&i18n.LocalizeConfig{
							DefaultMessage: &i18n.Message{
								ID:    "SignUp",
								Other: "Sign Up",
							},
						}),
						Style:    discordgo.PrimaryButton,
						CustomID: "boostRequest:signUp",
					},
					discordgo.Button{
						Label: m.localizer.MustLocalize(&i18n.LocalizeConfig{
							DefaultMessage: &i18n.Message{
								ID:    "Steal",
								Other: "Steal",
							},
						}),
						CustomID: "boostRequest:steal",
						Style:    discordgo.PrimaryButton,
					},
				},
			},
		},
	}, nil
}
