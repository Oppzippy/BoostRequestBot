package messages_test

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages/mocks"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func TestAdvertiserChosenDMToAdvertiserHuman(t *testing.T) {
	// TODO improve test
	m := messages.NewAdvertiserChosenDMToAdvertiser(
		emptyLocalizer(),
		&mocks.MockUserProvider{
			Value: &discordgo.User{
				ID:            "0987",
				Username:      "test",
				Discriminator: "1234",
			},
		},
		messages.NewDiscountFormatter(
			emptyLocalizer(),
			&mocks.MockRoleNameProvider{},
		),
		&repository.BoostRequest{
			RequesterID: "1",
			Channel:     repository.BoostRequestChannel{},
		},
	)
	_, err := m.Message()
	if err != nil {
		t.Errorf("error generating message: %v", err)
	}
}

func TestAdvertiserChosenDMToAdvertiserBot(t *testing.T) {
	// TODO improve test
	m := messages.NewAdvertiserChosenDMToAdvertiser(
		emptyLocalizer(),
		&mocks.MockUserProvider{
			Value: &discordgo.User{
				ID:            "0987",
				Username:      "test",
				Discriminator: "1234",
			},
		},
		messages.NewDiscountFormatter(
			emptyLocalizer(),
			&mocks.MockRoleNameProvider{},
		),
		&repository.BoostRequest{
			RequesterID: "1",
			EmbedFields: []*repository.MessageEmbedField{
				{
					Name:  "test",
					Value: "test",
				},
			},
			Channel: repository.BoostRequestChannel{},
		},
	)
	_, err := m.Message()
	if err != nil {
		t.Errorf("error generating message: %v", err)
	}
}
