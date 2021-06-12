package message_generator_test

import (
	"strings"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/oppzippy/BoostRequestBot/boost_request/message_generator"
	"github.com/oppzippy/BoostRequestBot/boost_request/message_generator/mocks"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func TestAdvertisreChosenDMToRequester(t *testing.T) {
	t.Parallel()
	br := &repository.BoostRequest{}
	m := message_generator.NewAdvertiserChosenDMToRequester(
		emptyLocalizer(),
		&mocks.MockDMUserProvider{
			Value: &discordgo.User{
				ID:            "1111",
				Username:      "test",
				Discriminator: "1234",
			},
		},
		message_generator.NewDiscountFormatter(
			emptyLocalizer(),
			&mocks.MockRoleNameProvider{},
		),
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
