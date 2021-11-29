package messages_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
	"github.com/oppzippy/BoostRequestBot/roll"
)

func TestBoostRequestRollMessage(t *testing.T) {
	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("generate uuid: %v", err)
		return
	}
	br := &repository.BoostRequest{
		Channel:    &repository.BoostRequestChannel{},
		ExternalID: &id,
	}
	roll := roll.NewWeightedRoll(3)
	roll.AddItem("advertiser1", 1)
	roll.AddItem("advertiser2", 2)
	roll.AddItem("advertiser3", 3)
	rollResults, ok := roll.Roll()
	if !ok {
		t.Error("roll wasn't ok")
		return
	}

	m := messages.NewBoostRequestRollMessage(
		emptyLocalizer(),
		br,
		rollResults,
	)

	message, err := m.Message()
	if err != nil {
		t.Errorf("error generating message: %v", err)
		return
	}

	lines := strings.Split(message.Embed.Description, "\n")
	var acc float64
	for i, line := range lines {
		if line != "" {
			indexFromOne := i + 1
			acc += float64(indexFromOne)
			expected := fmt.Sprintf(
				"<@advertiser%d>: %s to %s",
				indexFromOne,
				localizeFloat(acc-float64(indexFromOne)),
				localizeFloat(acc),
			)
			if strings.Index(line, expected) != 0 {
				t.Errorf("line is wrong: %s, expected %s", line, expected)
			}
		}
	}
}

func localizeFloat(f float64) string {
	return emptyLocalizer().MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Float64",
			Other: "{{.Float}}",
		},
		TemplateData: map[string]float64{
			"Float": f,
		},
	})
}
