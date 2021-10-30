package messenger

import (
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const copperPerSilver = 1_00
const copperPerGold = 1_00_00

func FormatCopper(localizer *i18n.Localizer, totalCopper int64) string {
	copper := totalCopper % 100
	silver := (totalCopper / copperPerSilver) % 100
	gold := totalCopper / copperPerGold

	parts := make([]string, 0, 3)

	if gold != 0 {
		parts = append(parts, localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "moneyFormatGold",
				Other: "{{.Gold}}g",
			},
			TemplateData: map[string]int64{
				"Gold": gold,
			},
		}))
	}
	if silver != 0 {
		parts = append(parts, localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "moneyFormatSilver",
				Other: "{{.Silver}}s",
			},
			TemplateData: map[string]int64{
				"Silver": silver,
			},
		}))
	}
	if copper != 0 {
		parts = append(parts, localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "moneyFormatCopper",
				Other: "{{.Copper}}c",
			},
			TemplateData: map[string]int64{
				"Copper": copper,
			},
		}))
	}

	return strings.Join(parts, " ")
}
