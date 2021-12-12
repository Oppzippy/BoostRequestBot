package messenger

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type MessageBroker struct {
	discord   DiscordSenderAndDeleter
	waitGroup *sync.WaitGroup
	quit      chan struct{}
	destroyed bool
}

func NewMessageBroker(discord DiscordSenderAndDeleter) *MessageBroker {
	return &MessageBroker{
		discord:   discord,
		waitGroup: new(sync.WaitGroup),
		quit:      make(chan struct{}),
	}
}

func (mb *MessageBroker) Send(dest *MessageDestination, mg MessageGenerator) (*discordgo.Message, error) {
	m := NewMessage(dest, mg)
	sentMessage, err := m.Send(mb.discord)
	return sentMessage, err
}

func (mb *MessageBroker) SendDelayed(
	dest *MessageDestination,
	mg MessageGenerator,
	delay time.Duration,
	cancel <-chan struct{},
) (<-chan *discordgo.Message, <-chan error) {
	m := NewMessage(dest, mg)
	dm := NewAsyncMessage(NewDelayedMessage(m, delay, cancel))

	return dm.Send(mb.discord)
}

func (mb *MessageBroker) SendTemporaryMessage(dest *MessageDestination, mg MessageGenerator, duration time.Duration) (*discordgo.Message, <-chan error) {
	errChannel := make(chan error, 1)
	m := NewMessage(dest, mg)
	sentMessage, err := m.Send(mb.discord)
	if err != nil {
		errChannel <- fmt.Errorf("sending temporary message: %v", err)
		close(errChannel)
		return nil, errChannel
	}
	mb.waitGroup.Add(1)
	go func() {
		defer mb.waitGroup.Done()
		defer close(errChannel)

		select {
		case <-time.After(duration):
		case <-mb.quit:
		}

		err := mb.discord.ChannelMessageDelete(sentMessage.ChannelID, sentMessage.ID)
		if err != nil {
			errChannel <- fmt.Errorf("error deleting temporary message: %v", err)
			return
		}
	}()
	return sentMessage, errChannel
}

func (mb *MessageBroker) Destroy() {
	if !mb.destroyed {
		mb.destroyed = true
		close(mb.quit)
		mb.waitGroup.Wait()
	}
}
