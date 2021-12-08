package application_commands

import "github.com/bwmarrin/discordgo"

var BoostRequestCommand = &discordgo.ApplicationCommand{
	Name:        "boostrequest",
	Description: "Boost Request Bot commands",
	Type:        discordgo.ChatApplicationCommand,
	Options: []*discordgo.ApplicationCommandOption{
		AutoSignupSubCommand,
	},
}
