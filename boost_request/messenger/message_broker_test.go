package messenger_test

import (
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/golang/mock/gomock"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages"
	"github.com/oppzippy/BoostRequestBot/boost_request/messenger"
	"github.com/oppzippy/BoostRequestBot/boost_request/messenger/mock_messenger"
)

func TestMessageBrokerSend(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	discord := mock_messenger.NewMockDiscordSenderAndDeleter(mockController)
	discord.
		EXPECT().
		ChannelMessageSendComplex("1", gomock.Any()).
		Return(&discordgo.Message{}, nil)

	mb := messenger.NewMessageBroker(discord)
	m, err := mb.Send(&messenger.MessageDestination{
		DestinationID:   "1",
		DestinationType: messenger.DestinationChannel,
	}, messages.NewStaticMessage(&discordgo.MessageSend{
		Content: "test",
	}))

	if err != nil {
		t.Errorf("error sending message: %v", err)
		return
	}
	if m == nil {
		t.Error("message is nil")
	}
}

func TestMessageBrokerSendDelayed(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	discord := mock_messenger.NewMockDiscordSenderAndDeleter(mockController)
	discord.
		EXPECT().
		ChannelMessageSendComplex("1", gomock.Any()).
		Return(&discordgo.Message{}, nil)

	mb := messenger.NewMessageBroker(discord)

	messageChannel, errChannel := mb.SendDelayed(
		&messenger.MessageDestination{
			DestinationID:   "1",
			DestinationType: messenger.DestinationChannel,
		},
		messages.NewStaticMessage(&discordgo.MessageSend{
			Content: "test",
		}),
		200*time.Millisecond,
		make(<-chan struct{}),
	)
	startTime := time.Now()

	select {
	case m, ok := <-messageChannel:
		if ok {
			elapsed := time.Since(startTime)
			if m == nil {
				t.Error("message is nil")
			} else if elapsed < 190*time.Millisecond {
				t.Errorf("expected 200ms to have passed, got %d", elapsed/time.Millisecond)
			} else if elapsed > 250*time.Millisecond {
				t.Errorf("expected 200ms to have passed, got %d", elapsed/time.Millisecond)
			}
			return
		}
	case err, ok := <-errChannel:
		if ok && err != nil {
			t.Errorf("error sending delayed message: %v", err)
			return
		}
	}
	t.Error("didn't get a value from either channel")
}

func TestMessageBrokerSendTemporary(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	discord := mock_messenger.NewMockDiscordSenderAndDeleter(mockController)
	discord.
		EXPECT().
		ChannelMessageSendComplex("1", gomock.Any()).
		Return(&discordgo.Message{
			ChannelID: "1",
			ID:        "2",
		}, nil)

	deleteChannel := make(chan struct{})
	discord.
		EXPECT().
		ChannelMessageDelete("1", "2").
		Return(nil).
		Do(func(channelID, messageID string) {
			deleteChannel <- struct{}{}
		})

	mb := messenger.NewMessageBroker(discord)
	m, errChannel := mb.SendTemporaryMessage(&messenger.MessageDestination{
		DestinationID:   "1",
		DestinationType: messenger.DestinationChannel,
	}, messages.NewStaticMessage(&discordgo.MessageSend{
		Content: "test",
	}), 200*time.Millisecond)
	startTime := time.Now()

	select {
	case err := <-errChannel:
		t.Errorf("error sending message: %v", err)
		return
	default:
	}
	if m == nil {
		t.Error("message is nil")
		return
	}

	select {
	case <-deleteChannel:
		elapsed := time.Since(startTime)
		if elapsed < 190*time.Millisecond {
			t.Errorf("expected 200ms to have passed, got %v", elapsed/time.Millisecond)
		}
	case <-time.After(250 * time.Millisecond):
		t.Error("250ms passed without the message being deleted")
	}
}
