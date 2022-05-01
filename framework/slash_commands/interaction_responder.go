package slash_commands

import "github.com/bwmarrin/discordgo"

type interactionResponder interface {
	InteractionRespond(interaction *discordgo.Interaction, resp *discordgo.InteractionResponse) error
}
