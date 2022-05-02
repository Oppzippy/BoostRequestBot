package admin_commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/oppzippy/BoostRequestBot/util/pointers"
)

var RollChannelCommand = &discordgo.ApplicationCommand{
	Name:                     "rollchannel",
	Description:              "Boost request log channel administration.",
	Type:                     discordgo.ChatApplicationCommand,
	DefaultMemberPermissions: pointers.To(int64(discordgo.PermissionAdministrator)),
	DMPermission:             pointers.To(false),
	Options: []*discordgo.ApplicationCommandOption{
		setRollChannelSubCommand,
		removeRollChannelSubCommand,
	},
}

var setRollChannelSubCommand = &discordgo.ApplicationCommandOption{
	Name:        "set",
	Description: "Sets a channel to log RNG rolls to.",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	ChannelTypes: []discordgo.ChannelType{
		discordgo.ChannelTypeGuildCategory,
	},
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "channel",
			Description: "Roll channel.",
			Type:        discordgo.ApplicationCommandOptionChannel,
			Required:    true,
		},
	},
}

var removeRollChannelSubCommand = &discordgo.ApplicationCommandOption{
	Name:        "remove",
	Description: "Stops logging RNG rolls.",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	ChannelTypes: []discordgo.ChannelType{
		discordgo.ChannelTypeGuildCategory,
	},
}
