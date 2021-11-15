package messages

import (
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/shopspring/decimal"
)

const (
	copperPerSilver int64 = 1_00
	copperPerGold   int64 = 1_00_00
)

func formatCopper(localizer *i18n.Localizer, totalCopper int64) string {
	if totalCopper < copperPerGold*1000 {
		return formatCopperToGoldSilverCopper(localizer, totalCopper)
	} else {
		return formatCopperToThousandsOfGold(totalCopper)
	}
}

func formatCopperToGoldSilverCopper(localizer *i18n.Localizer, totalCopper int64) string {
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

func formatCopperToThousandsOfGold(copper int64) string {
	copperDecimal := decimal.NewFromInt(copper)
	thousandsOfGold := copperDecimal.Div(decimal.NewFromInt(copperPerGold * 1000))
	// TODO make configurable
	return thousandsOfGold.Round(2).String() + "k <:gold:909618212717592607>"
}
