//go:generate stringer -output=name_visibility_setting_string.go -type=NameVisibilitySetting -linecomment
package repository

type NameVisibilitySetting int

const (
	NameVisibilityShow          = NameVisibilitySetting(iota) // SHOW
	NameVisibilityShowInDMsOnly                               // SHOW_IN_DMS_ONLY
	NameVisibilityHide                                        // HIDE
)

func NameVisibilitySettingFromString(s string) NameVisibilitySetting {
	switch s {
	case NameVisibilityShow.String():
		return NameVisibilityShow
	case NameVisibilityShowInDMsOnly.String():
		return NameVisibilityShowInDMsOnly
	case NameVisibilityHide.String():
		return NameVisibilityHide
	}
	return NameVisibilityShow
}
