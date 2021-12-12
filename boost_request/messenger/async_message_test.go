package messenger_test

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/golang/mock/gomock"
	"github.com/oppzippy/BoostRequestBot/boost_request/messenger"
	"github.com/oppzippy/BoostRequestBot/boost_request/messenger/mock_messenger"
)

func TestAsyncMessageSend(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	discord := mock_messenger.NewMockDiscordSender(mockController)
	sendable := mock_messenger.NewMockSendable(mockController)

	sendable.
		EXPECT().
		Send(discord).
		Return(&discordgo.Message{}, nil)

	m := messenger.NewAsyncMessage(sendable)
	messageChannel, errChannel := m.Send(discord)

	var err error
	var message *discordgo.Message
	select {
	case message = <-messageChannel:
	case err = <-errChannel:
	}
	if err != nil {
		t.Errorf("failed to send: %v", err)
		return
	}
	if message == nil {
		t.Error("message is nil")
	}
}
