package message_generator_test

type MockRoleNameProvider struct {
	MockRoleName string
}

func (rnp *MockRoleNameProvider) RoleName(guildID, roleID string) string {
	return rnp.MockRoleName
}
