package messages

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type CreditsUpdatedDM struct {
	localizer *i18n.Localizer
	credits   int
}

func NewCreditsUpdatedDM(
	localizer *i18n.Localizer, credits int,
) *CreditsUpdatedDM {
	return &CreditsUpdatedDM{
		localizer: localizer,
		credits:   credits,
	}
}

func (m *CreditsUpdatedDM) Message() (*discordgo.MessageSend, error) {
	return &discordgo.MessageSend{
		Content: m.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "YourCreditsUpdated",
				One:   "You now have {{.Credits}} boost request steal credit.",
				Other: "You now have {{.Credits}} boost request steal credit.",
			},
			TemplateData: map[string]int{
				"Credits": m.credits,
			},
			PluralCount: m.credits,
		}),
	}, nil
}
