package messenger

import "github.com/bwmarrin/discordgo"

type AsyncMessage struct {
	sendable Sendable
}

func NewAsyncMessage(sendable Sendable) *AsyncMessage {
	return &AsyncMessage{
		sendable: sendable,
	}
}

func (m *AsyncMessage) Send(discord DiscordSender) (<-chan *discordgo.Message, <-chan error) {
	messageChannel := make(chan *discordgo.Message, 1)
	errChannel := make(chan error, 1)
	go func() {
		defer close(messageChannel)
		defer close(errChannel)

		sentMessage, err := m.sendable.Send(discord)
		if err != nil {
			errChannel <- err
		}
		if sentMessage != nil {
			messageChannel <- sentMessage
		}
	}()
	return messageChannel, errChannel
}
