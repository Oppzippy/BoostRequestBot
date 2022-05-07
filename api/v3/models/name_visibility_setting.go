//go:generate stringer -output=name_visibility_setting_string.go -type=NameVisibilitySetting -linecomment
package models

import (
	"encoding/json"
)

type NameVisibilitySetting int

const (
	NameVisibilityShowInDMsOnly = NameVisibilitySetting(iota) // SHOW_IN_DMS_ONLY
	NameVisibilityShow                                        // SHOW
	NameVisibilityHide                                        // HIDE
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
	case NameVisibilityShow.String():
		return NameVisibilityShow
	case NameVisibilityShowInDMsOnly.String():
		return NameVisibilityShowInDMsOnly
	case NameVisibilityHide.String():
		return NameVisibilityHide
	}
	var def NameVisibilitySetting
	return def
}
