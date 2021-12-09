package messages

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type AutoSignUpExpiredMessage struct {
	localizer *i18n.Localizer
}

func NewAutoSignUpExpiredMessage(
	localizer *i18n.Localizer,
) *AutoSignUpExpiredMessage {
	return &AutoSignUpExpiredMessage{
		localizer: localizer,
	}
}

func (m *AutoSignUpExpiredMessage) Message() (*discordgo.MessageSend, error) {
	return &discordgo.MessageSend{
		Content: m.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "AutoSignUpExpired",
				Other: "Auto sign up expired.",
			},
		}),
	}, nil
}
