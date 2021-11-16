package messages_test

import (
	"strings"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages/mocks"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages/partials"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func TestBackendAdvertiserChosenMessage(t *testing.T) {
	t.Parallel()
	br := &repository.BoostRequest{
		AdvertiserID: "123",
		Message:      "boost please!",
		Channel:      repository.BoostRequestChannel{},
	}
	m := messages.NewBackendAdvertiserChosenMessage(
		emptyLocalizer(),
		&mocks.MockUserProvider{
			Value: &discordgo.User{},
		},
		partials.NewDiscountFormatter(
			emptyLocalizer(),
			&mocks.MockRoleNameProvider{
				Value: "",
			},
		),
		br,
	)
	message, err := m.Message()
	if err != nil {
		t.Errorf("error generating message: %v", err)
		return
	}
	if message.Embed.Thumbnail.URL == "" {
		t.Errorf("thumbnail was not set")
	}
	if !strings.Contains(message.Embed.Description, "123") {
		t.Errorf("the chosen advertiser is not mentioned in the message: %s", message.Embed.Description)
	}
}
