package admin_commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/oppzippy/BoostRequestBot/util/pointers"
)

var LogChannelCommand = &discordgo.ApplicationCommand{
	Name:                     "logchannel",
	Description:              "Boost request log channel administration.",
	Type:                     discordgo.ChatApplicationCommand,
	DefaultMemberPermissions: pointers.To(int64(discordgo.PermissionAdministrator)),
	DMPermission:             pointers.To(false),
	Options: []*discordgo.ApplicationCommandOption{
		setLogChannelSubCommand,
		removeLogChannelSubCommand,
	},
}

var setLogChannelSubCommand = &discordgo.ApplicationCommandOption{
	Name:        "set",
	Description: "All boost requests will be logged to this channel upon creation.",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	ChannelTypes: []discordgo.ChannelType{
		discordgo.ChannelTypeGuildCategory,
	},
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "channel",
			Description: "Log channel.",
			Type:        discordgo.ApplicationCommandOptionChannel,
			Required:    true,
		},
	},
}

var removeLogChannelSubCommand = &discordgo.ApplicationCommandOption{
	Name:        "remove",
	Description: "Stops logging.",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	ChannelTypes: []discordgo.ChannelType{
		discordgo.ChannelTypeGuildCategory,
	},
}
