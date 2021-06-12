package mocks

import "github.com/bwmarrin/discordgo"

type MockUserProvider struct {
	Value *discordgo.User
}

func (m *MockUserProvider) User(userID string) (*discordgo.User, error) {
	return m.Value, nil
}
