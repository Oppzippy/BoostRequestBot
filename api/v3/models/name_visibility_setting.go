//go:generate stringer -output=name_visibility_setting_string.go -type=NameVisibilitySetting -linecomment
package models

import (
	"encoding/json"
)

type NameVisibilitySetting int

const (
	ShowInDMsOnly = NameVisibilitySetting(iota) // SHOW_IN_DMS_ONLY
	Show                                        // SHOW
	Hide                                        // HIDE
)

func (s NameVisibilitySetting) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *NameVisibilitySetting) UnmarshalJSON(b []byte) error {
	var setting string
	err := json.Unmarshal(b, &setting)
	if err != nil {
		return err
	}
	value := NameVisibilitySettingFromString(setting)
	*s = value
	return nil
}

func NameVisibilitySettingFromString(s string) NameVisibilitySetting {
	switch s {
	case Show.String():
		return Show
	case ShowInDMsOnly.String():
		return ShowInDMsOnly
	case Hide.String():
		return Hide
	}
	var def NameVisibilitySetting
	return def
}
