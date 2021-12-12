//go:generate mockgen -source discord.go -destination mock_messenger/discord_mock.go
package messenger

import "github.com/bwmarrin/discordgo"

type DiscordSender interface {
	discordChannelMessageSendComplex
	discordUserChannelCreate
}

type DiscordSenderAndDeleter interface {
	DiscordSender
	discordChannelMessageDelete
}

type discordChannelMessageSendComplex interface {
	ChannelMessageSendComplex(channelID string, data *discordgo.MessageSend) (st *discordgo.Message, err error)
}

type discordUserChannelCreate interface {
	UserChannelCreate(recipientID string) (st *discordgo.Channel, err error)
}

type discordChannelMessageDelete interface {
	ChannelMessageDelete(channelID, messageID string) (err error)
}
