package partials

type roleNameProvider interface {
	RoleName(guildID, roleID string) string
}
