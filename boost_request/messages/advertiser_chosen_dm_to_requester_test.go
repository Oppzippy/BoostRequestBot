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

func TestAdvertisreChosenDMToRequester(t *testing.T) {
	t.Parallel()
	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("generate uuid: %v", err)
		return
	}
	br := &repository.BoostRequest{
		ExternalID: &id,
	}
	m := messages.NewAdvertiserChosenDMToRequester(
		emptyLocalizer(),
		&mocks.MockUserProvider{
			Value: &discordgo.User{
				ID:            "1111",
				Username:      "test",
				Discriminator: "1234",
			},
		},
		br,
	)

	message, err := m.Message()
	if err != nil {
		t.Errorf("error generating message: %v", err)
		return
	}
	if !strings.Contains(message.Embed.Description, "test#1234") {
		t.Errorf("message should contain the advertiser's tag in case discord hasn't cached them: %v", message.Embed.Description)
	}
	if !strings.Contains(message.Embed.Description, "<@1111>") {
		t.Errorf("the advertiser should be mentioned in the message: %v", message.Embed.Description)
	}
	if message.Embed.Thumbnail.URL == "" {
		t.Errorf("the thumbnail should be set to the advertiser's avatar")
	}
}
