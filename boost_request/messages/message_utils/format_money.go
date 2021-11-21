package message_utils

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/shopspring/decimal"
)

const (
	copperPerSilver int64 = 1_00
	copperPerGold   int64 = 1_00_00
)

func FormatCopper(localizer *i18n.Localizer, copper int64) string {
	copperDecimal := decimal.NewFromInt(copper)
	thousandsOfGold := copperDecimal.Div(decimal.NewFromInt(copperPerGold * 1000))
	// TODO make emoji configurable
	return thousandsOfGold.Round(2).String() + "k <:gold:909618212717592607>"
}
