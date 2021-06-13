package messages_test

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages/mocks"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func TestLogChannelMessage(t *testing.T) {
	br := &repository.BoostRequest{
		Channel: repository.BoostRequestChannel{},
	}
	m := messages.NewLogChannelMessage(emptyLocalizer(), &mocks.MockUserProvider{
		Value: &discordgo.User{},
	}, br)

	_, err := m.Message()
	if err != nil {
		t.Errorf("error generating message: %v", err)
	}

	// TODO add more tests
}
