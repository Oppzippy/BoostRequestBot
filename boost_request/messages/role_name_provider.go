package messages

type RoleNameProvider interface {
	RoleName(guildID, roleID string) string
}
