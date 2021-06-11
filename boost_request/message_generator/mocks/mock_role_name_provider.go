package mocks

type MockRoleNameProvider struct {
	Value string
}

func (rnp *MockRoleNameProvider) RoleName(guildID, roleID string) string {
	return rnp.Value
}
