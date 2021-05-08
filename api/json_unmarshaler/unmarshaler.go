package json_unmarshaler

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator"
)

type structValidator interface {
	Struct(interface{}) error
}

type Unmarshaler struct {
	Validate structValidator
}

func New() *Unmarshaler {
	return &Unmarshaler{
		Validate: validator.New(),
	}
}

// Unmarshal will unmarshal and validate the passed json data.
// If validation fails, dest will still have been modified!
func (p *Unmarshaler) Unmarshal(data []byte, dest interface{}) error {
	err := json.Unmarshal(data, dest)
	if err != nil {
		return err
	}
	err = p.Validate.Struct(dest)
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
	return err
}
