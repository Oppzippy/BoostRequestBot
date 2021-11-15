package boost_request

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type InteractionRegistry struct {
	discord          *discordgo.Session
	handlers         []interactionHandler
	destoryFunctions []func()
}

type interactionHandler interface {
	Matches(discord *discordgo.Session, event *discordgo.InteractionCreate) bool
	Handle(discord *discordgo.Session, event *discordgo.InteractionCreate) error
}

func NewInteractionRegistry(discord *discordgo.Session) *InteractionRegistry {
	r := &InteractionRegistry{
		discord:          discord,
		destoryFunctions: make([]func(), 0),
	}

	r.destoryFunctions = append(r.destoryFunctions, discord.AddHandler(r.onInteractionCreate))
	return r
}

func (r *InteractionRegistry) AddHandler(handler interactionHandler) {
	r.handlers = append(r.handlers, handler)
}

func (r *InteractionRegistry) RemoveAllHandlers() {
	r.handlers = make([]interactionHandler, 0)
}

func (r *InteractionRegistry) Destroy() {
	for _, f := range r.destoryFunctions {
		f()
	}
}

func (r *InteractionRegistry) onInteractionCreate(discord *discordgo.Session, event *discordgo.InteractionCreate) {
	for _, handler := range r.handlers {
		if handler.Matches(discord, event) {
			err := handler.Handle(discord, event)
			if err != nil {
				log.Printf("error handling interaction: %v", err)
			}
		}
	}
}
