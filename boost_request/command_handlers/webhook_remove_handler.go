package command_handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type WebhookRemoveHandler struct {
	bundle *i18n.Bundle
	repo   repository.Repository
}

func NewWebhookRemoveHandler(bundle *i18n.Bundle, repo repository.Repository) *WebhookRemoveHandler {
	return &WebhookRemoveHandler{
		bundle: bundle,
		repo:   repo,
	}
}

func (h *WebhookRemoveHandler) Handle(event *discordgo.InteractionCreate, options map[string]*discordgo.ApplicationCommandInteractionDataOption) (*discordgo.InteractionResponse, error) {
	localizer := i18n.NewLocalizer(h.bundle, string(event.Locale))
	webhook, err := h.repo.GetWebhook(event.GuildID)
	if err != nil {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "WebhookNotSet",
						Other: "A webhook is not set.",
					},
				}),
				Flags: uint64(discordgo.MessageFlagsEphemeral),
			},
		}, nil
	}
	err = h.repo.DeleteWebhook(webhook)
	if err != nil {
		return nil, fmt.Errorf("error deleting webhook: %w", err)
	}
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "WebhookRemoved",
					Other: "Webhook removed.",
				},
			}),
			Flags: uint64(discordgo.MessageFlagsEphemeral),
		},
	}, nil
}
