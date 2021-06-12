package message_generator_test

import (
	"strings"
	"testing"

	"github.com/oppzippy/BoostRequestBot/boost_request/message_generator"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func TestBackendAdvertiserChosenMessage(t *testing.T) {
	t.Parallel()
	br := &repository.BoostRequest{
		AdvertiserID: "123",
		Message:      "boost please!",
		Channel:      repository.BoostRequestChannel{},
	}
	m := message_generator.NewBackendAdvertiserChosenMessage(emptyLocalizer(), br, "https://www.example.com")
	message, err := m.Message()
	if err != nil {
		t.Errorf("error generating message: %v", err)
		return
	}
	if message.Embed.Thumbnail.URL != "https://www.example.com" {
		t.Errorf("thumbnail was not set")
	}
	if !strings.Contains(message.Embed.Description, "123") {
		t.Errorf("the chosen advertiser is not mentioned in the message: %s", message.Embed.Description)
	}
}
