package messages

import (
	"github.com/bwmarrin/discordgo"
)

type StaticMessage struct {
	message *discordgo.MessageSend
}

func NewStaticMessage(message *discordgo.MessageSend) *StaticMessage {
	return &StaticMessage{
		message: message,
	}
}

func (m *StaticMessage) Message() (*discordgo.MessageSend, error) {
	return m.message, nil
}
