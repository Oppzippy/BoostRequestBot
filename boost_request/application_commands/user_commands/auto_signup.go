package user_commands

import "github.com/bwmarrin/discordgo"

var autoSignupSubCommand = &discordgo.ApplicationCommandOption{
	Name:        "autosignup",
	Description: "Automatically sign up for all new boost requests for a limited period of time.",
	Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
	ChannelTypes: []discordgo.ChannelType{
		discordgo.ChannelTypeGuildCategory,
	},
	Options: []*discordgo.ApplicationCommandOption{
		autoSignupStartSubCommand,
		autoSignupCancelSubCommand,
	},
}

var autoSignupStartSubCommand = &discordgo.ApplicationCommandOption{
	Name:        "start",
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

var autoSignupCancelSubCommand = &discordgo.ApplicationCommandOption{
	Name:        "stop",
	Description: "Stop automatically signing up for all new boost requests.",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	ChannelTypes: []discordgo.ChannelType{
		discordgo.ChannelTypeGuildCategory,
	},
}
