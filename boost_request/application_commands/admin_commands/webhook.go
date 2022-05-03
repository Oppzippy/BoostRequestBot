package admin_commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/oppzippy/BoostRequestBot/util/pointers"
)

var WebhookCommand = &discordgo.ApplicationCommand{
	Name:                     "webhook",
	Description:              "A specified URL will receive a POST request every time a webhook event occurs.",
	Type:                     discordgo.ChatApplicationCommand,
	DefaultMemberPermissions: pointers.To(int64(discordgo.PermissionAdministrator)),
	DMPermission:             pointers.To(false),
	DefaultPermission:        pointers.To(false),
	Options: []*discordgo.ApplicationCommandOption{
		listWebhookSubCommand,
		setWebhookSubCommand,
		removeWebhookSubCommand,
	},
}

var listWebhookSubCommand = &discordgo.ApplicationCommandOption{
	Name:        "list",
	Description: "Lists the webhook url that you set.",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	ChannelTypes: []discordgo.ChannelType{
		discordgo.ChannelTypeGuildCategory,
	},
}

var setWebhookSubCommand = &discordgo.ApplicationCommandOption{
	Name:        "set",
	Description: "Sets a URL that will receive a POST request every time a webhook event occurs.",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	ChannelTypes: []discordgo.ChannelType{
		discordgo.ChannelTypeGuildCategory,
	},
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "webhook-url",
			Description: "Webhook URL.",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
		},
	},
}

var removeWebhookSubCommand = &discordgo.ApplicationCommandOption{
	Name:        "remove",
	Description: "Deactivates webhooks.",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	ChannelTypes: []discordgo.ChannelType{
		discordgo.ChannelTypeGuildCategory,
	},
}
