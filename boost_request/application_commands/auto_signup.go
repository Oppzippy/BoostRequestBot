package application_commands

import "github.com/bwmarrin/discordgo"

var AutoSignupSubCommand = &discordgo.ApplicationCommandOption{
	Name:        "autosignup",
	Description: "Automatically sign up for all new boost requests for a limited period of time.",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	ChannelTypes: []discordgo.ChannelType{
		discordgo.ChannelTypeGuildCategory,
	},
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "duration",
			Description: "Duration in minutes",
			Type:        discordgo.ApplicationCommandOptionInteger,
			Required:    false,
		},
	},
}
