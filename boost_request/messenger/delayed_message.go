package messenger

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type DelayedMessage struct {
	message Sendable
	delay   time.Duration
	cancel  <-chan struct{}
}

func NewDelayedMessage(sendable Sendable, delay time.Duration, cancel <-chan struct{}) *DelayedMessage {
	return &DelayedMessage{
		message: sendable,
		delay:   delay,
		cancel:  cancel,
	}
}

func (m *DelayedMessage) Send(discord DiscordSender) (*discordgo.Message, error) {
	select {
	case <-time.After(m.delay):
		return m.message.Send(discord)
	case <-m.cancel:
		return nil, nil
	}
}
