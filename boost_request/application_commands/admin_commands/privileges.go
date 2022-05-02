package admin_commands

import (
	"math"

	"github.com/bwmarrin/discordgo"
	"github.com/oppzippy/BoostRequestBot/util/pointers"
)

var PrivilegesCommand = &discordgo.ApplicationCommand{
	Name:                     "privileges",
	Description:              "Administer boost request privileges.",
	Type:                     discordgo.ChatApplicationCommand,
	DefaultMemberPermissions: pointers.To(int64(discordgo.PermissionAdministrator)),
	DMPermission:             pointers.To(false),
	Options: []*discordgo.ApplicationCommandOption{
		listPrivilegesSubCommand,
		removePrivilegesSubCommand,
		setPrivilegesSubCommand,
	},
}

var listPrivilegesSubCommand = &discordgo.ApplicationCommandOption{
	Name:        "list",
	Description: "Lists all privileges.",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	ChannelTypes: []discordgo.ChannelType{
		discordgo.ChannelTypeGuildCategory,
	},
}

var removePrivilegesSubCommand = &discordgo.ApplicationCommandOption{
	Name:        "remove",
	Description: "Removes a role's boost request privileges.",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	ChannelTypes: []discordgo.ChannelType{
		discordgo.ChannelTypeGuildCategory,
	},
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "role",
			Description: "Role to remove all privileges from.",
			Type:        discordgo.ApplicationCommandOptionRole,
			Required:    true,
		},
	},
}

var setPrivilegesSubCommand = &discordgo.ApplicationCommandOption{
	Name:        "set",
	Description: "Sets the privileges for an advertiser role.",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	ChannelTypes: []discordgo.ChannelType{
		discordgo.ChannelTypeGuildCategory,
	},
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "role",
			Description: "Role to grant privileges to.",
			Type:        discordgo.ApplicationCommandOptionRole,
			Required:    true,
		},
		{
			Name:        "weight",
			Description: "When randomly choosing an advertiser, the weight will be applied to advertisers of this role.",
			Type:        discordgo.ApplicationCommandOptionInteger,
			MinValue:    pointers.To(math.SmallestNonzeroFloat64),
			Required:    true,
		},
		{
			Name:        "delay-in-seconds",
			Description: "Delay in seconds after the creation of a boost request before it can be claimed.",
			Type:        discordgo.ApplicationCommandOptionInteger,
			MinValue:    pointers.To(float64(1)),
			Required:    true,
		},
		{
			Name:        "auto-signup-duration-in-minutes",
			Description: "Maximum number of minutes that advertisers can enable auto signup for.",
			Type:        discordgo.ApplicationCommandOptionInteger,
			MinValue:    pointers.To(float64(1)),
			Required:    false,
		},
	},
}
