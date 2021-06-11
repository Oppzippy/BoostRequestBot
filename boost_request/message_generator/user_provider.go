package message_generator

import "github.com/bwmarrin/discordgo"

type userProvider interface {
	User(userID string) (*discordgo.User, error)
}
