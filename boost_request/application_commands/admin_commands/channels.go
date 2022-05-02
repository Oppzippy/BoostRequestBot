package admin_commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/oppzippy/BoostRequestBot/util/pointers"
)

var ChannelsCommand = &discordgo.ApplicationCommand{
	Name:                     "channels",
	Description:              "Boost request channel administration.",
	Type:                     discordgo.ChatApplicationCommand,
	DefaultMemberPermissions: pointers.To(int64(discordgo.PermissionAdministrator)),
	DMPermission:             pointers.To(false),
	Options: []*discordgo.ApplicationCommandOption{
		listChannelsSubCommand,
		addChannelSubCommand,
		removeChannelSubCommand,
	},
}

var listChannelsSubCommand = &discordgo.ApplicationCommandOption{
	Name:        "list",
	Description: "Lists all boost request channels on the server.",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	ChannelTypes: []discordgo.ChannelType{
		discordgo.ChannelTypeGuildCategory,
	},
}

var addChannelSubCommand = &discordgo.ApplicationCommandOption{
	Name:        "add",
	Description: "Adds a new boost request channel.",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	ChannelTypes: []discordgo.ChannelType{
		discordgo.ChannelTypeGuildCategory,
	},
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type: discordgo.ApplicationCommandOptionChannel,
			Name: "frontend-channel",
			NameLocalizations: map[discordgo.Locale]string{
				"en-US": "frontend-channel",
			},
			Description: "Boost requests will be requested in this channel.",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionChannel,
			Name:        "backend-channel",
			Description: "Boost requests will be posted to this channel for advertisers to claim.",
			Required:    true,
		},
	},
}

var removeChannelSubCommand = &discordgo.ApplicationCommandOption{
	Name:        "remove",
	Description: "Deletes a boost request channel.",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	ChannelTypes: []discordgo.ChannelType{
		discordgo.ChannelTypeGuildCategory,
	},
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionChannel,
			Name:        "frontend-channel",
			Description: "Boost requests are requested in this channel.",
			Required:    true,
		},
	},
}
