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
	boostRequest        *repository.BoostRequest
	localizer           *i18n.Localizer
	embedPartial        *partials.BoostRequestEmbedTemplate
	buttonConfiguration BackendSignupMessageButtonConfiguration
}

type BackendSignupMessageButtonConfiguration struct {
	SignUp       bool
	Steal        bool
	CancelSignup bool
}

func NewBackendSignupMessage(
	localizer *i18n.Localizer,
	br *repository.BoostRequest,
	buttonConfiguration BackendSignupMessageButtonConfiguration,
) *BackendSignupMessage {
	return &BackendSignupMessage{
		boostRequest:        br,
		localizer:           localizer,
		embedPartial:        partials.NewBoostRequestEmbedTemplate(localizer, br),
		buttonConfiguration: buttonConfiguration,
	}
}

func (m *BackendSignupMessage) Message() (*discordgo.MessageSend, error) {
	br := m.boostRequest
	embed, err := m.embedPartial.Embed(partials.BoostRequestEmbedConfiguration{
		Price: true,
		ID:    true,
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
		mentions := make([]string, 0, len(br.PreferredAdvertiserIDs))
		for id := range br.PreferredAdvertiserIDs {
			mentions = append(mentions, fmt.Sprintf("<@%s>", id))
		}
		title := m.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "PreferredClaimer",
				One:   "Preferred Claimer",
				Other: "Preferred Claimers",
			},
			PluralCount: len(br.PreferredAdvertiserIDs),
		})
		preferredAdvertiserMentions = fmt.Sprintf("**%s:** %s", title, strings.Join(mentions, " "))
	}

	components := make([]discordgo.MessageComponent, 0, 5)
	if m.buttonConfiguration.SignUp {
		components = append(components, discordgo.Button{
			Label: m.localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "SignUp",
					Other: "Sign Up",
				},
			}),
			Style:    discordgo.PrimaryButton,
			CustomID: "boostRequest:signUp",
		})
	}
	if m.buttonConfiguration.Steal {
		components = append(components, discordgo.Button{
			Label: m.localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "Steal",
					Other: "Steal",
				},
			}),
			CustomID: "boostRequest:steal",
			Style:    discordgo.PrimaryButton,
		})
	}
	if m.buttonConfiguration.CancelSignup {
		components = append(components, discordgo.Button{
			Label: m.localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "CancelSignup",
					Other: "Cancel Signup",
				},
			}),
			CustomID: "boostRequest:cancelSignup",
			Style:    discordgo.SecondaryButton,
		})
	}

	return &discordgo.MessageSend{
		Content: preferredAdvertiserMentions,
		Embed:   embed,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: components,
			},
		},
	}, nil
}
