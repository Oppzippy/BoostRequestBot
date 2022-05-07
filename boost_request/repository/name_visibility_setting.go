//go:generate stringer -output=name_visibility_setting_string.go -type=NameVisibilitySetting -linecomment
package repository

type NameVisibilitySetting int

const (
	Show          = NameVisibilitySetting(iota) // SHOW
	ShowInDMsOnly                               // SHOW_IN_DMS_ONLY
	Hide                                        // HIDE
)

func NameVisibilitySettingFromString(s string) NameVisibilitySetting {
	switch s {
	case Show.String():
		return Show
	case ShowInDMsOnly.String():
		return ShowInDMsOnly
	case Hide.String():
		return Hide
	}
	return Show
}
