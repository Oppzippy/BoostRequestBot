package slash_commands

import "github.com/bwmarrin/discordgo"

type SlashCommandMiddleware func(
	interaction *discordgo.InteractionCreate,
	options map[string]*discordgo.ApplicationCommandInteractionDataOption,
) (*discordgo.InteractionResponse, error)

type SlashCommandHandler = SlashCommandMiddleware

func buildPipeline(middleware []SlashCommandMiddleware, commandHandler SlashCommandHandler) []SlashCommandMiddleware {
	pipeline := make([]SlashCommandMiddleware, 0, len(middleware)+1)
	pipeline = append(pipeline, middleware...)
	if commandHandler != nil {
		pipeline = append(pipeline, commandHandler)
	}
	return pipeline
}
