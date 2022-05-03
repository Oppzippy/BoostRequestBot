package command_handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type LogChannelRemoveHandler struct {
	bundle *i18n.Bundle
	repo   repository.Repository
}

func NewLogChannelRemoveHandler(bundle *i18n.Bundle, repo repository.Repository) *LogChannelRemoveHandler {
	return &LogChannelRemoveHandler{
		bundle: bundle,
		repo:   repo,
	}
}

func (h *LogChannelRemoveHandler) Handle(event *discordgo.InteractionCreate, options map[string]*discordgo.ApplicationCommandInteractionDataOption) (*discordgo.InteractionResponse, error) {
	localizer := i18n.NewLocalizer(h.bundle, string(event.Locale))
	err := h.repo.DeleteLogChannel(event.GuildID)
	if err != nil {
		return nil, fmt.Errorf("error deleting boost request log channel: %w", err)
	}
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "BoostRequestLoggingDisabled",
					Other: "Boost request logging is disabled.",
				},
			}),
			Flags: uint64(discordgo.MessageFlagsEphemeral),
		},
	}, nil
}
