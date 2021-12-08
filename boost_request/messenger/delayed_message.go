package messenger

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type sendable interface {
	Send(discord *discordgo.Session) (*discordgo.Message, error)
}

type delayedMessage struct {
	message sendable
	delay   time.Duration
}

func newDelayedMessage(sendable sendable, delay time.Duration) *delayedMessage {
	return &delayedMessage{
		message: sendable,
		delay:   delay,
	}
}

func (m *delayedMessage) Send(discord *discordgo.Session) (*discordgo.Message, error) {
	time.Sleep(m.delay)
	return m.message.Send(discord)
}
