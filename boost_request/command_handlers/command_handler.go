package command_handlers

import "github.com/bwmarrin/discordgo"

type CommandHandler interface {
	Handle(event *discordgo.InteractionCreate, options map[string]*discordgo.ApplicationCommandInteractionDataOption) (*discordgo.InteractionResponse, error)
}
