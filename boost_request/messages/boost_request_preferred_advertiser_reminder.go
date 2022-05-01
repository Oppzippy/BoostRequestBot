package messages

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages/partials"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type BoostRequestPreferredAdvertiserReminder struct {
	localizer    *i18n.Localizer
	boostRequest *repository.BoostRequest
	embedPartial *partials.BoostRequestEmbedTemplate
}

func NewBoostRequestPreferredAdvertiserReminder(
	localizer *i18n.Localizer, br *repository.BoostRequest,
) *BoostRequestPreferredAdvertiserReminder {
	return &BoostRequestPreferredAdvertiserReminder{
		localizer:    localizer,
		boostRequest: br,
		embedPartial: partials.NewBoostRequestEmbedTemplate(localizer, br),
	}
}

func (m *BoostRequestPreferredAdvertiserReminder) Message() (*discordgo.MessageSend, error) {
	return &discordgo.MessageSend{
		Content: m.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "PreferredAdvertiserAFK",
				Other: "You have a preferred advertiser set and 15 minutes have passed without the boost request being claimed. If you wish to remove your advertiser preference for this request, please use the Remove Advertiser Preference button on the message above.",
			},
		}),
	}, nil
}
