package messages_test

import (
	"strings"
	"testing"

	"github.com/oppzippy/BoostRequestBot/boost_request/messages"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages/mocks"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
	"github.com/shopspring/decimal"
)

func TestBackendSignupMessage(t *testing.T) {
	t.Parallel()
	br := &repository.BoostRequest{
		Message: "Boost please",
		Channel: repository.BoostRequestChannel{
			BackendChannelID: "1",
		},
	}

	bsm := messages.NewBackendSignupMessage(
		emptyLocalizer(),
		&messages.DiscountFormatter{},
		br,
	)

	t.Run("Message", func(t *testing.T) {
		message, err := bsm.Message()
		if err != nil {
			t.Errorf("error generating message: %v", err)
		} else if message.Embed.Description != br.Message {
			t.Errorf("expected description %s, got %s", br.Message, message.Embed.Description)
		}
	})
}

func TestBackendSignupMessageRoleDiscount(t *testing.T) {
	t.Parallel()
	discount, err := decimal.NewFromString("0.2")
	if err != nil {
		t.Errorf("parsing discount: %v", err)
	}

	br := &repository.BoostRequest{
		Channel: repository.BoostRequestChannel{},
		RoleDiscounts: []*repository.RoleDiscount{
			{
				RoleID:    "1",
				BoostType: "mythic+",
				Discount:  discount,
			},
			{
				RoleID:    "1",
				BoostType: "raid",
				Discount:  discount,
			},
		},
	}
	bsm := messages.NewBackendSignupMessage(
		emptyLocalizer(),
		messages.NewDiscountFormatter(
			emptyLocalizer(),
			&mocks.MockRoleNameProvider{
				Value: "booster",
			},
		),
		br,
	)

	message, err := bsm.Message()
	if err != nil {
		t.Errorf("error generating message: %v", err)
		return
	}
	lines := strings.Split(message.Embed.Fields[0].Value, "\n")

	if expected := "20% discount on mythic+ (booster)"; lines[0] != expected {
		t.Errorf("Expected %s, got %s", expected, lines[0])
	}
	if expected := "20% discount on raid (booster)"; lines[1] != expected {
		t.Errorf("Expected %s, got %s", expected, lines[1])
	}
}
