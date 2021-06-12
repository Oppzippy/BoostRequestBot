package message_generator_test

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/oppzippy/BoostRequestBot/boost_request/message_generator"
	"github.com/oppzippy/BoostRequestBot/boost_request/message_generator/mocks"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func TestAdvertiserChosenDMToAdvertiserHuman(t *testing.T) {
	// TODO improve test
	m := message_generator.NewAdvertiserChosenDMToAdvertiser(
		emptyLocalizer(),
		&mocks.MockUserProvider{
			Value: &discordgo.User{
				ID:            "0987",
				Username:      "test",
				Discriminator: "1234",
			},
		},
		message_generator.NewDiscountFormatter(
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
	m := message_generator.NewAdvertiserChosenDMToAdvertiser(
		emptyLocalizer(),
		&mocks.MockUserProvider{
			Value: &discordgo.User{
				ID:            "0987",
				Username:      "test",
				Discriminator: "1234",
			},
		},
		message_generator.NewDiscountFormatter(
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
