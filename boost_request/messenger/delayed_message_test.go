package messenger_test

import (
	"errors"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/golang/mock/gomock"
	"github.com/oppzippy/BoostRequestBot/boost_request/messenger"
	"github.com/oppzippy/BoostRequestBot/boost_request/messenger/mock_messenger"
)

func TestDelayedMessageSendFuture(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	discord := mock_messenger.NewMockDiscordSender(mockController)
	sendable := mock_messenger.NewMockSendable(mockController)

	sendable.EXPECT().Send(discord).Return(&discordgo.Message{}, nil)

	dm := messenger.NewDelayedMessage(sendable, 100*time.Millisecond, make(<-chan struct{}))

	result := make(chan error)
	go func() {
		message, err := dm.Send(discord)
		if err != nil {
			result <- err
			return
		}
		if message == nil {
			result <- errors.New("message is nil")
			return
		}
		close(result)
	}()
	select {
	case err, ok := <-result:
		if ok {
			t.Errorf("error sending message: %v", err)
		}
	case <-time.After(200 * time.Millisecond):
		t.Error("delayed message timed out")
	}
}
