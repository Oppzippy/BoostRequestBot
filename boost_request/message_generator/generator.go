package message_generator

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type Generator struct {
	localizer *i18n.Localizer
	discord   *discordgo.Session
}

func NewGenerator(localizer *i18n.Localizer, discord *discordgo.Session) *Generator {
	return &Generator{
		localizer: localizer,
		discord:   discord,
	}
}

func (gen *Generator) BackendSignupMessage(br *repository.BoostRequest) *BackendSignupMessage {
	return NewBackendSignupMessage(gen.localizer, gen.discountFormatter(), br)
}

func (gen *Generator) BoostRequestCreatedDM(br *repository.BoostRequest) *BoostRequestCreatedDM {
	return NewBoostRequestCreatedDM(gen.localizer, gen.userProvider(), br)
}

func (gen *Generator) DMBlockedMessage(userID string) *DMBlockedMessage {
	return NewDMBlockedMessage(gen.localizer, userID)
}

func (gen *Generator) discountFormatter() *DiscountFormatter {
	return NewDiscountFormatter(gen.localizer, gen.roleNameProvider())
}

func (gen *Generator) roleNameProvider() RoleNameProvider {
	return NewDiscordRoleNameProvider(gen.discord)
}

func (gen *Generator) userProvider() userProvider {
	return gen.discord
}
