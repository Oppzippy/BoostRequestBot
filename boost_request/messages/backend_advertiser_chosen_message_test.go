package messages_test

import (
	"strings"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages/mocks"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func TestBackendAdvertiserChosenMessage(t *testing.T) {
	t.Parallel()
	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("generate uuid: %v", err)
		return
	}
	br := &repository.BoostRequest{
		AdvertiserID: "123",
		Message:      "boost please!",
		Channel:      &repository.BoostRequestChannel{},
		ExternalID:   &id,
	}
	m := messages.NewBackendAdvertiserChosenMessage(
		emptyLocalizer(),
		&mocks.MockUserProvider{
			Value: &discordgo.User{},
		},
		br,
	)
	message, err := m.Message()
	if err != nil {
		t.Errorf("error generating message: %v", err)
		return
	}
	if message.Embed.Thumbnail != nil {
		t.Errorf("thumbnail was set but should not be to protect privacy")
	}
	if strings.Contains(message.Embed.Description, "123") {
		t.Errorf("the chosen advertiser is mentioned in the message but shoulnd't be: %s", message.Embed.Description)
	}
}
