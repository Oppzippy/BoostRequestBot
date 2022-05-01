package messenger_test

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/golang/mock/gomock"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages"
	"github.com/oppzippy/BoostRequestBot/boost_request/messenger"
	"github.com/oppzippy/BoostRequestBot/boost_request/messenger/mock_messenger"
)

func TestMessageSend(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	discord := mock_messenger.NewMockDiscordSender(mockController)

	m := messenger.NewMessage(&messenger.MessageDestination{
		DestinationID:   "1",
		DestinationType: messenger.DestinationChannel,
	}, messages.NewStaticMessage(&discordgo.MessageSend{
		Content: "test",
	}))
	discord.
		EXPECT().
		ChannelMessageSendComplex("1", gomock.Any()).
		Return(&discordgo.Message{}, nil)
	sentMessage, err := m.Send(discord)
	if err != nil {
		t.Errorf("send failed: %s", err)
		return
	}
	if sentMessage == nil {
		t.Error("returned message was nil")
		return
	}
}
