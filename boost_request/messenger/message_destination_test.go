package messenger_test

import (
	"errors"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/golang/mock/gomock"
	"github.com/oppzippy/BoostRequestBot/boost_request/messenger"
	"github.com/oppzippy/BoostRequestBot/boost_request/messenger/mock_messenger"
)

func TestMessageDestinationResolveChannel(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	discord := mock_messenger.NewMockdiscordUserChannelCreate(mockController)

	dest := messenger.MessageDestination{
		DestinationID:   "1",
		DestinationType: messenger.DestinationChannel,
	}
	channelID, err := dest.ResolveChannelID(discord)
	if err != nil {
		t.Errorf("failed to resovle channel id: %v", err)
		return
	}
	if channelID != "1" {
		t.Errorf("expected channel id 1, got %s", channelID)
	}
}

func TestMessageDestinationResolveDM(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	discord := mock_messenger.NewMockdiscordUserChannelCreate(mockController)
	discord.
		EXPECT().
		UserChannelCreate("1").
		Return(&discordgo.Channel{
			ID: "2",
		}, nil)

	dest := messenger.MessageDestination{
		DestinationID:   "1",
		DestinationType: messenger.DestinationUser,
	}
	channelID, err := dest.ResolveChannelID(discord)
	if err != nil {
		t.Errorf("failed to resolve channel id: %v", err)
		return
	}
	if channelID != "2" {
		t.Errorf("expected channel id 2, got %s", channelID)
	}
}

func TestMessageDestinationResolveFallbackDM(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	discord := mock_messenger.NewMockdiscordUserChannelCreate(mockController)
	discord.
		EXPECT().
		UserChannelCreate("1").
		Return(nil, errors.New("user not found"))

	dest := messenger.MessageDestination{
		DestinationID:     "1",
		DestinationType:   messenger.DestinationUser,
		FallbackChannelID: "2",
	}
	_, err := dest.ResolveChannelID(discord)
	if err == nil {
		t.Error("expected an error but got nil")
		return
	}
}
