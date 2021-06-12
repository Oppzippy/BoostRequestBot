package message_generator_test

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/oppzippy/BoostRequestBot/boost_request/message_generator"
	"github.com/oppzippy/BoostRequestBot/boost_request/message_generator/mocks"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func TestBoostRequestCreatedDM(t *testing.T) {
	t.Parallel()
	br := &repository.BoostRequest{
		Message: "Boost please",
		Channel: repository.BoostRequestChannel{
			SkipsBuyerDM: false,
		},
	}
	createdDM := message_generator.NewBoostRequestCreatedDM(
		emptyLocalizer(),
		&mocks.MockDMUserProvider{
			Value: &discordgo.User{
				Username:      "test",
				Discriminator: "1234",
			},
		},
		br,
	)
	m, _ := createdDM.Message()

	if m.Embed.Author.Name != "test#1234" {
		t.Errorf("expected tag test#1234, got %s", m.Embed.Author.Name)
	}
	if m.Embed.Description != "Boost please" {
		t.Errorf("expected description %s, got %s", br.Message, m.Embed.Description)
	}
}
