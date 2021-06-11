package message_generator

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type DMBlockedMessage struct {
	localizer *i18n.Localizer
	userID    string
}

func NewDMBlockedMessage(localizer *i18n.Localizer, userID string) *DMBlockedMessage {
	return &DMBlockedMessage{
		localizer: localizer,
		userID:    userID,
	}
}

func (m *DMBlockedMessage) Message() (*discordgo.MessageSend, error) {
	return &discordgo.MessageSend{
		Content: m.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				Other: "<@{{.UserID}}>, I can't DM you. Please allow DMs from server members by right clicking the server and enabling \"Allow direct messages from server members.\" in Privacy Settings.",
			},
			TemplateData: map[string]string{
				"UserID": m.userID,
			},
		}),
	}, nil
}
