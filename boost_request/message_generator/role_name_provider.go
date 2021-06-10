package message_generator

type RoleNameProvider interface {
	RoleName(guildID, roleID string) string
}
