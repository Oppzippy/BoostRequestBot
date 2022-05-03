package command_handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type RollChannelRemoveHandler struct {
	bundle *i18n.Bundle
	repo   repository.Repository
}

func NewRollChannelRemoveHandler(bundle *i18n.Bundle, repo repository.Repository) *RollChannelRemoveHandler {
	return &RollChannelRemoveHandler{
		bundle: bundle,
		repo:   repo,
	}
}

func (h *RollChannelRemoveHandler) Handle(event *discordgo.InteractionCreate, options map[string]*discordgo.ApplicationCommandInteractionDataOption) (*discordgo.InteractionResponse, error) {
	localizer := i18n.NewLocalizer(h.bundle, string(event.Locale))
	err := h.repo.DeleteRollChannel(event.GuildID)
	if err != nil {
		return nil, fmt.Errorf("error deleting boost request roll channel: %v", err)
	}
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "BoostRequestRollChannelRemoved",
					Other: "The boost request roll channel has been removed.",
				},
			}),
			Flags: uint64(discordgo.MessageFlagsEphemeral),
		},
	}, nil
}
