package message_generator

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type Generator struct {
	localiser *i18n.Localizer
	discord   *discordgo.Session
}

func NewGenerator(localizer *i18n.Localizer, discord *discordgo.Session) *Generator {
	return &Generator{
		localiser: localizer,
		discord:   discord,
	}
}

func (gen *Generator) BackendSignupMessage(br *repository.BoostRequest) *BackendSignupMessage {
	return NewBackendSignupMessage(gen.localiser, gen.discountFormatter(), br)
}

func (gen *Generator) discountFormatter() *DiscountFormatter {
	return NewDiscountFormatter(gen.localiser, gen.roleNameProvider())
}

func (gen *Generator) roleNameProvider() RoleNameProvider {
	return NewDiscordRoleNameProvider(gen.discord)
}
