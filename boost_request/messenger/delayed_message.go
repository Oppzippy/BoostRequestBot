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
	cancel  <-chan struct{}
}

func newDelayedMessage(sendable sendable, delay time.Duration, cancel <-chan struct{}) *delayedMessage {
	return &delayedMessage{
		message: sendable,
		delay:   delay,
		cancel:  cancel,
	}
}

func (m *delayedMessage) Send(discord *discordgo.Session) (*discordgo.Message, error) {
	select {
	case <-time.After(m.delay):
		return m.message.Send(discord)
	case <-m.cancel:
		return nil, nil
	}
}
