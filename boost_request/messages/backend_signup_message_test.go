package messages_test

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func TestBackendSignupMessage(t *testing.T) {
	t.Parallel()
	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("generate uuid: %v", err)
		return
	}
	br := &repository.BoostRequest{
		Message: "Boost please",
		Channel: &repository.BoostRequestChannel{
			BackendChannelID: "1",
		},
		ExternalID: &id,
	}

	bsm := messages.NewBackendSignupMessage(
		emptyLocalizer(),
		br,
		messages.BackendSignupMessageButtonConfiguration{
			SignUp:       true,
			Steal:        true,
			CancelSignup: true,
		},
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
	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("generate uuid: %v", err)
		return
	}
	if err != nil {
		t.Errorf("parsing discount: %v", err)
	}

	br := &repository.BoostRequest{
		Channel:    &repository.BoostRequestChannel{},
		ExternalID: &id,
	}
	bsm := messages.NewBackendSignupMessage(
		emptyLocalizer(),
		br,
		messages.BackendSignupMessageButtonConfiguration{
			SignUp:       true,
			Steal:        true,
			CancelSignup: true,
		},
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
