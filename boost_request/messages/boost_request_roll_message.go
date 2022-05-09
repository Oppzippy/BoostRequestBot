package messages

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
	"github.com/oppzippy/BoostRequestBot/util/weighted_picker"
)

type BoostRequestRollMessage struct {
	localizer    *i18n.Localizer
	boostRequest *repository.BoostRequest
	rollResults  *weighted_picker.WeightedPickerResults[string]
}

func NewBoostRequestRollMessage(
	localizer *i18n.Localizer, br *repository.BoostRequest, rollResults *weighted_picker.WeightedPickerResults[string],
) *BoostRequestRollMessage {
	return &BoostRequestRollMessage{
		localizer:    localizer,
		boostRequest: br,
		rollResults:  rollResults,
	}
}

func (m *BoostRequestRollMessage) Message() (*discordgo.MessageSend, error) {
	if m.rollResults == nil {
		return nil, errors.New("rollResults must not be nil")
	}

	sb := strings.Builder{}
	var weightAccumulator float64
	for iter := m.rollResults.Iterator(); iter.HasNext(); {
		advertiserID, weight, isChosenItem := iter.Next()
		weightAccumulator += weight

		sb.WriteString(m.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "AdvertiserRollRange",
				Description: "Min is inclusive, Max is exclusive.",
				Other:       "{{.AdvertiserMention}}: {{.Min}} to {{.Max}}",
			},
			TemplateData: map[string]interface{}{
				"AdvertiserMention": fmt.Sprintf("<@%s>", advertiserID),
				"Min":               weightAccumulator - weight,
				"Max":               weightAccumulator,
			},
		}))
		if isChosenItem {
			sb.WriteString(m.localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "AdvertiserRollRangeChosen",
					Other: "   **<-- {{.Roll}}**",
				},
				TemplateData: map[string]float64{
					"Roll": m.rollResults.ChosenNumber(),
				},
			}))
		}
		sb.WriteString("\n")
	}

	return &discordgo.MessageSend{
		Content: m.boostRequest.Message,
		Embed: &discordgo.MessageEmbed{
			Title: m.localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "RollResults",
					Other: "Roll Results",
				},
			}),
			Description: sb.String(),
		},
		AllowedMentions: &discordgo.MessageAllowedMentions{},
	}, nil
}
