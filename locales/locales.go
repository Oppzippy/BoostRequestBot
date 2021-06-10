package locales

import (
	"embed"
	"fmt"
	"io"
	"log"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

//go:embed *.yaml
var locales embed.FS

func Bundle() *i18n.Bundle {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	err := addLocalesToBundle(bundle)
	if err != nil {
		log.Fatalf("failed to register locales: %v", err)
	}

	return bundle
}

func addLocalesToBundle(bundle *i18n.Bundle) error {
	dirEntries, err := locales.ReadDir(".")
	if err != nil {
		return err
	}
	for _, dirEntry := range dirEntries {
		file, err := locales.Open(dirEntry.Name())
		if err != nil {
			return fmt.Errorf("opening %s: %v", dirEntry.Name(), err)
		}
		content, err := io.ReadAll(file)
		if err != nil {
			return fmt.Errorf("reading %s: %v", dirEntry.Name(), err)
		}
		_, err = bundle.ParseMessageFileBytes(content, dirEntry.Name())
		if err != nil {
			return fmt.Errorf("parsing locale file %s: %v", dirEntry.Name(), err)
		}
	}
	return nil
}
