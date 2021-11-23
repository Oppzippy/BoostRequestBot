package boost_request

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type InteractionRegistry struct {
	discord          *discordgo.Session
	bundle           *i18n.Bundle
	handlers         []interactionHandler
	destoryFunctions []func()
}

type interactionHandler interface {
	Matches(discord *discordgo.Session, event *discordgo.InteractionCreate) bool
	Handle(discord *discordgo.Session, event *discordgo.InteractionCreate, localizer *i18n.Localizer) error
}

func NewInteractionRegistry(discord *discordgo.Session, bundle *i18n.Bundle) *InteractionRegistry {
	r := &InteractionRegistry{
		discord:          discord,
		bundle:           bundle,
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
			localizer := i18n.NewLocalizer(r.bundle, "en")
			err := handler.Handle(discord, event, localizer)
			if err != nil {
				log.Printf("error handling interaction: %v", err)
			}
		}
	}
}
