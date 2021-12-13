package messages

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type AutoSignupExpiredMessage struct {
	localizer *i18n.Localizer
}

func NewAutoSignupExpiredMessage(
	localizer *i18n.Localizer,
) *AutoSignupExpiredMessage {
	return &AutoSignupExpiredMessage{
		localizer: localizer,
	}
}

func (m *AutoSignupExpiredMessage) Message() (*discordgo.MessageSend, error) {
	return &discordgo.MessageSend{
		Content: m.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "AutoSignupExpired",
				Other: "Auto sign up expired.",
			},
		}),
	}, nil
}
