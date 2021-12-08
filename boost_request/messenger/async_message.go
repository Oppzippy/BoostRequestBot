package messenger

import "github.com/bwmarrin/discordgo"

type asyncMessage struct {
	sendable sendable
}

func newAsyncMessage(sendable sendable) *asyncMessage {
	return &asyncMessage{
		sendable: sendable,
	}
}

func (m *asyncMessage) Send(discord *discordgo.Session) (<-chan *discordgo.Message, <-chan error) {
	messageChannel := make(chan *discordgo.Message, 1)
	errChannel := make(chan error, 1)
	go func() {
		defer close(messageChannel)
		defer close(errChannel)

		sentMessage, err := m.sendable.Send(discord)
		if err != nil {
			errChannel <- err
			return
		}
		messageChannel <- sentMessage
	}()
	return messageChannel, errChannel
}
