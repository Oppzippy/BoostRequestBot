package mocks

import "github.com/bwmarrin/discordgo"

type MockDMUserProvider struct {
	Value *discordgo.User
}

func (m *MockDMUserProvider) User(userID string) (*discordgo.User, error) {
	return m.Value, nil
}
