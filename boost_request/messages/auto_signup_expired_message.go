package messages

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type AutoSignupExpiredMessage struct {
	localizer *i18n.Localizer
	guildID   string
}

func NewAutoSignupExpiredMessage(
	localizer *i18n.Localizer,
	guildID string,
) *AutoSignupExpiredMessage {
	return &AutoSignupExpiredMessage{
		localizer: localizer,
		guildID:   guildID,
	}
}

func (m *AutoSignupExpiredMessage) Message() (*discordgo.MessageSend, error) {
	return &discordgo.MessageSend{
		Content: m.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "AutoSignupExpired",
				Other: "Auto signup expired.",
			},
		}),
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Enable Auto Signup Again",
						Style:    discordgo.PrimaryButton,
						CustomID: fmt.Sprintf("autoSignup:%s", m.guildID),
					},
				},
			},
		},
	}, nil
}
