//go:generate mockgen -source sendable.go -destination mock_messenger/sendable.go -aux_files messenger=discord.go
package messenger

import "github.com/bwmarrin/discordgo"

type Sendable interface {
	Send(discord DiscordSender) (*discordgo.Message, error)
}
