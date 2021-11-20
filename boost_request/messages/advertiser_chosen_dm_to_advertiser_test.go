package messages_test

import (
	"strings"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages/mocks"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages/partials"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func TestAdvertiserChosenDMToAdvertiserHuman(t *testing.T) {
	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("generate uuid: %v", err)
		return
	}
	m := messages.NewAdvertiserChosenDMToAdvertiser(
		emptyLocalizer(),
		&mocks.MockUserProvider{
			Value: &discordgo.User{
				ID:            "0987",
				Username:      "test",
				Discriminator: "1234",
			},
		},
		partials.NewDiscountFormatter(
			emptyLocalizer(),
			&mocks.MockRoleNameProvider{},
		),
		&repository.BoostRequest{
			RequesterID: "1",
			Channel:     repository.BoostRequestChannel{},
			ExternalID:  &id,
		},
	)
	message, err := m.Message()
	if err != nil {
		t.Errorf("error generating message: %v", err)
	}

	if !strings.Contains(message.Embed.Description, "test#1234") {
		t.Errorf("buyer's name not found in message: %s", message.Embed.Description)
	}
	if !strings.Contains(message.Embed.Description, "<@0987>") {
		t.Errorf("buyer mention not found in message: %s", message.Embed.Description)
	}
}

func TestAdvertiserChosenDMToAdvertiserBot(t *testing.T) {
	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("generate uuid: %v", err)
		return
	}
	m := messages.NewAdvertiserChosenDMToAdvertiser(
		emptyLocalizer(),
		&mocks.MockUserProvider{
			Value: &discordgo.User{
				ID:            "0987",
				Username:      "test",
				Discriminator: "1234",
			},
		},
		partials.NewDiscountFormatter(
			emptyLocalizer(),
			&mocks.MockRoleNameProvider{},
		),
		&repository.BoostRequest{
			RequesterID: "1",
			EmbedFields: []*repository.MessageEmbedField{
				{
					Name:  "Buyer",
					Value: "Test",
				},
			},
			Channel:    repository.BoostRequestChannel{},
			ExternalID: &id,
		},
	)
	message, err := m.Message()
	if err != nil {
		t.Errorf("error generating message: %v", err)
	}

	if field := message.Embed.Fields[0]; field.Name != "Buyer" || field.Value != "Test" {
		t.Errorf("fields don't match the boost request: %s=%s", field.Name, field.Value)
	}
}
