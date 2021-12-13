package messages

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type AutoSignupExpiringSoonMessage struct {
	localizer *i18n.Localizer
	timeLeft  time.Duration
}

func NewAutoSignupExpiringSoonMessage(
	localizer *i18n.Localizer,
	timeLeft time.Duration,
) *AutoSignupExpiringSoonMessage {
	return &AutoSignupExpiringSoonMessage{
		localizer: localizer,
		timeLeft:  timeLeft,
	}
}

func (m *AutoSignupExpiringSoonMessage) Message() (*discordgo.MessageSend, error) {
	return &discordgo.MessageSend{
		Content: m.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "AutoSignUpExpiringSoon",
				Other: "In {{ .TimeLeft }}, you will no longer automatically sign up for boost requests.",
			},
			TemplateData: map[string]interface{}{
				"TimeLeft": m.timeLeft,
			},
		}),
	}, nil
}
