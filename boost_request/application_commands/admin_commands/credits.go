package admin_commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/oppzippy/BoostRequestBot/util/pointers"
)

var CreditsCommand = &discordgo.ApplicationCommand{
	Name:                     "credits",
	Description:              "Administer boost request steal credits.",
	Type:                     discordgo.ChatApplicationCommand,
	DefaultMemberPermissions: pointers.To(int64(discordgo.PermissionAdministrator)),
	DMPermission:             pointers.To(false),
	DefaultPermission:        pointers.To(false),
	Options: []*discordgo.ApplicationCommandOption{
		addCreditsSubCommand,
		setCreditsSubCommand,
		checkCreditsSubCommand,
	},
}

var addCreditsSubCommand = &discordgo.ApplicationCommandOption{
	Name:        "add",
	Description: "Adds boost request steal credits to a user.",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	ChannelTypes: []discordgo.ChannelType{
		discordgo.ChannelTypeGuildCategory,
	},
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "user",
			Description: "User who should be granted credits.",
			Type:        discordgo.ApplicationCommandOptionUser,
			Required:    true,
		},
		{
			Name:        "credits",
			Description: "Number of credits to add. Enter a negative number to subtract credits.",
			Type:        discordgo.ApplicationCommandOptionInteger,
			Required:    true,
		},
	},
}

var setCreditsSubCommand = &discordgo.ApplicationCommandOption{
	Name:        "set",
	Description: "Sets the number of boost request steal credits available to a user.",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	ChannelTypes: []discordgo.ChannelType{
		discordgo.ChannelTypeGuildCategory,
	},
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "user",
			Description: "User to change the number of available credits credits of.",
			Type:        discordgo.ApplicationCommandOptionUser,
			Required:    true,
		},
		{
			Name:        "credits",
			Description: "New number of credits.",
			Type:        discordgo.ApplicationCommandOptionInteger,
			MinValue:    pointers.To(float64(0)),
			Required:    true,
		},
	},
}

var checkCreditsSubCommand = &discordgo.ApplicationCommandOption{
	Name:        "list",
	Description: "Lists the number of boost request steal credits available to a user.",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	ChannelTypes: []discordgo.ChannelType{
		discordgo.ChannelTypeGuildCategory,
	},
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "user",
			Description: "User to check the number of available credits credits of.",
			Type:        discordgo.ApplicationCommandOptionUser,
			Required:    true,
		},
	},
}
