package slash_commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type SlashCommandRegistry struct {
	Logger          *log.Logger
	commandHandlers map[string]SlashCommandHandler
	middleware      []SlashCommandMiddleware
}

func NewSlashCommandRegistry() *SlashCommandRegistry {
	registry := &SlashCommandRegistry{
		commandHandlers: map[string]SlashCommandHandler{},
		middleware:      []SlashCommandMiddleware{},
		Logger:          log.Default(),
	}

	return registry
}

func (r *SlashCommandRegistry) AttachToDiscord(discord *discordgo.Session) (removeHandler func()) {
	return discord.AddHandler(r.OnInteraction)
}

func (r *SlashCommandRegistry) RegisterCommand(path string, handler SlashCommandHandler) {
	r.commandHandlers[path] = handler
}

func (r *SlashCommandRegistry) RegisterMidleware(middleware SlashCommandMiddleware) {
	r.middleware = append(r.middleware, middleware)
}

func (r *SlashCommandRegistry) OnInteraction(responder interactionResponder, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	commandData := i.ApplicationCommandData()
	commandPath := CommandPathToString(&commandData)
	commandHandler := r.commandHandlers[commandPath]
	if commandHandler == nil {
		return
	}

	pipeline := buildPipeline(r.middleware, commandHandler)
	response := r.runPipeline(pipeline, i, &commandData)
	if response == nil {
		r.Logger.Printf("slash command didn't return a response: %v", commandPath)
		response.Data = &discordgo.InteractionResponseData{
			Content: "An error has occurred.",
		}
	}

	err := responder.InteractionRespond(i.Interaction, response)
	if err != nil {
		r.Logger.Printf("error running slash command %s: sending response: %v", commandPath, err)
	}
}

func (r *SlashCommandRegistry) runPipeline(
	pipeline []SlashCommandMiddleware,
	interaction *discordgo.InteractionCreate,
	commandData *discordgo.ApplicationCommandInteractionData,
) *discordgo.InteractionResponse {
	options := parseOptions(commandData)

	for i, middleware := range r.middleware {
		response, err := middleware(interaction, options)
		if err != nil {
			r.Logger.Printf("error running slash command %s: pipeline step %d: %v", CommandPathToString(commandData), i, err)
		}
		if response != nil {
			return response
		}
	}
	return nil
}
