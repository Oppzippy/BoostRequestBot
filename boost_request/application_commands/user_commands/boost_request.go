package user_commands

import "github.com/bwmarrin/discordgo"

var dmPermission = false

var BoostRequestCommand = &discordgo.ApplicationCommand{
	Name:         "boostrequest",
	Description:  "Boost Request Bot commands",
	Type:         discordgo.ChatApplicationCommand,
	DMPermission: &dmPermission,
	Options: []*discordgo.ApplicationCommandOption{
		autoSignupSubCommand,
	},
}
