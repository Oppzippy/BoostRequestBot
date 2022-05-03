package command_handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type StealCreditsSetHandler struct {
	bundle *i18n.Bundle
	repo   repository.Repository
}

func NewCreditsSetHandler(bundle *i18n.Bundle, repo repository.Repository) *StealCreditsSetHandler {
	return &StealCreditsSetHandler{
		bundle: bundle,
		repo:   repo,
	}
}

func (h *StealCreditsSetHandler) Handle(event *discordgo.InteractionCreate, options map[string]*discordgo.ApplicationCommandInteractionDataOption) (*discordgo.InteractionResponse, error) {
	localizer := i18n.NewLocalizer(h.bundle, string(event.Locale))
	user := options["user"].UserValue(nil)
	amount := int(options["credits"].IntValue())

	err := h.repo.UpdateStealCreditsForUser(event.GuildID, user.ID, amount)
	if err != nil {
		return nil, fmt.Errorf("error updating steal credits: %w", err)
	}
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "StealCreditsSet",
					Other: "Set steal credits to {.Credits}",
				},
				PluralCount: amount,
			}),
			Flags: uint64(discordgo.MessageFlagsEphemeral),
		},
	}, nil
}
