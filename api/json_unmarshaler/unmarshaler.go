package json_unmarshaler

import (
	"encoding/json"
	"io"
	"log"

	en2 "github.com/go-playground/validator/v10/translations/en"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"

	"github.com/go-playground/validator/v10"
)

type Unmarshaler struct {
	Validate   *validator.Validate
	translator ut.Translator
}

func New() *Unmarshaler {
	unmarshaler := &Unmarshaler{
		Validate: validator.New(),
	}

	defaultLang := en.New()
	uni := ut.New(defaultLang, defaultLang)
	trans, _ := uni.GetTranslator("en")
	unmarshaler.translator = trans
	err := en2.RegisterDefaultTranslations(unmarshaler.Validate, trans)
	if err != nil {
		log.Printf("failed to register translations: %v", err)
	}

	return unmarshaler
}

// Unmarshal will unmarshal and validate the passed json data.
// If validation fails, dest will still have been modified!
func (p *Unmarshaler) Unmarshal(data []byte, dest interface{}) error {
	err := json.Unmarshal(data, dest)
	if err != nil {
		return err
	}
	err = p.Validate.Struct(dest)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		err = p.translateValidationErrors(validationErrors)
	}
	return err
}

// UnmarshalReader will unmarshal and validate the passed json data.
// If validation fails, dest will still have been modified!
func (p *Unmarshaler) UnmarshalReader(r io.ReadCloser, dest interface{}) error {
	err := json.NewDecoder(r).Decode(dest)
	if err != nil {
		return err
	}
	err = p.Validate.Struct(dest)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		err = p.translateValidationErrors(validationErrors)
	}

	return err
}

func (p *Unmarshaler) translateValidationErrors(validationErrors validator.ValidationErrors) error {
	translatedValidationErrors := validationErrors.Translate(p.translator)
	return ValidationError{validationErrors.Error(), translatedValidationErrors}
}

type ValidationError struct {
	message          string
	TranslatedErrors map[string]string
}

func (err ValidationError) Error() string {
	return err.message
}
