package messages_test

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages/mocks"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages/partials"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func TestBoostRequestCreatedDM(t *testing.T) {
	t.Parallel()
	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("generate uuid: %v", err)
		return
	}
	br := &repository.BoostRequest{
		Message: "Boost please",
		Channel: &repository.BoostRequestChannel{
			SkipsBuyerDM: false,
		},
		ExternalID: &id,
	}
	createdDM := messages.NewBoostRequestCreatedDM(
		emptyLocalizer(),
		&mocks.MockUserProvider{
			Value: &discordgo.User{
				Username:      "test",
				Discriminator: "1234",
			},
		},
		partials.NewDiscountFormatter(emptyLocalizer(), &mocks.MockRoleNameProvider{
			Value: "",
		}),
		br,
	)
	m, _ := createdDM.Message()

	if m.Embed.Author.Name != "test#1234" {
		t.Errorf("expected tag test#1234, got %s", m.Embed.Author.Name)
	}
	if m.Embed.Fields[0].Value != "Boost please" {
		t.Errorf("expected description %s, got %s", br.Message, m.Embed.Description)
	}
}
