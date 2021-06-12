package messages_test

import (
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
		},
	}

	t.Run("NoRole", func(t *testing.T) {
		bsm := messages.NewBackendSignupMessage(
			emptyLocalizer(),
			messages.NewDiscountFormatter(
				emptyLocalizer(),
				&mocks.MockRoleNameProvider{},
			),
			br,
		)

		message, err := bsm.Message()
		if err != nil {
			t.Errorf("error generating message: %v", err)
			return
		}
		actual := message.Embed.Fields[0].Value

		if expected := "20% discount on mythic+"; actual != expected {
			t.Errorf("Expected %s, got %s", expected, actual)
		}
	})

	t.Run("WithRole", func(t *testing.T) {
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
		actual := message.Embed.Fields[0].Value

		if expected := "20% discount on mythic+ (booster)"; actual != expected {
			t.Errorf("Expected %s, got %s", expected, actual)
		}
	})
}
