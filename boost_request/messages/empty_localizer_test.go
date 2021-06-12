package messages_test

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func emptyLocalizer() *i18n.Localizer {
	return i18n.NewLocalizer(i18n.NewBundle(language.AmericanEnglish), "en")
}
