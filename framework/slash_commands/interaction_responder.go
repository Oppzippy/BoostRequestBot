//go:generate mockgen -source interaction_responder.go -destination mock_slash_commands/interaction_responder.go
package slash_commands

import "github.com/bwmarrin/discordgo"

type InteractionResponder interface {
	InteractionRespond(interaction *discordgo.Interaction, resp *discordgo.InteractionResponse) error
}
