package command_handlers

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type CreditsListHandler struct {
	bundle *i18n.Bundle
	repo   repository.Repository
}

func NewCreditsListHandler(bundle *i18n.Bundle, repo repository.Repository) *CreditsListHandler {
	return &CreditsListHandler{
		bundle: bundle,
		repo:   repo,
	}
}

func (h *CreditsListHandler) Handle(event *discordgo.InteractionCreate, options map[string]*discordgo.ApplicationCommandInteractionDataOption) (*discordgo.InteractionResponse, error) {
	localizer := i18n.NewLocalizer(h.bundle, string(event.Locale))

	user := options["user"].UserValue(nil)

	credits, err := h.repo.GetStealCreditsForUser(event.GuildID, user.ID)
	if errors.Is(err, repository.ErrNoResults) {
		credits = 0
	} else if err != nil {
		return nil, fmt.Errorf("error fetching boost request steal credits for user in admin check credits command: %w", err)
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "UserStealCredits",
					One:   "{{.User}} has {{.Credits}} steal credit.",
					Other: "{{.User}} has {{.Credits}} steal credits.",
				},
				TemplateData: map[string]interface{}{
					"User":    user.Mention(),
					"Credits": credits,
				},
				PluralCount: credits,
			}),
			Flags: uint64(discordgo.MessageFlagsEphemeral),
		},
	}, nil
}
