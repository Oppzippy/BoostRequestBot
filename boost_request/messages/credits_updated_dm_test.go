package messages_test

import (
	"strings"
	"testing"

	"github.com/oppzippy/BoostRequestBot/boost_request/messages"
)

func TestCreditsUpdatedDM(t *testing.T) {
	m := messages.NewCreditsUpdatedDM(emptyLocalizer(), 3)
	message, err := m.Message()
	if err != nil {
		t.Errorf("error generating message: %v", err)
		return
	}

	if !strings.Contains(message.Content, "3") {
		t.Errorf("number of credits was not found in the message: %v", message.Content)
	}
}
