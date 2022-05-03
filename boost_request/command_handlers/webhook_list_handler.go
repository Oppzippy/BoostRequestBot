package command_handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type WebhookListHandler struct {
	bundle *i18n.Bundle
	repo   repository.Repository
}

func NewWebhookListHandler(bundle *i18n.Bundle, repo repository.Repository) *WebhookListHandler {
	return &WebhookListHandler{
		bundle: bundle,
		repo:   repo,
	}
}

func (h *WebhookListHandler) Handle(event *discordgo.InteractionCreate, options map[string]*discordgo.ApplicationCommandInteractionDataOption) (*discordgo.InteractionResponse, error) {
	localizer := i18n.NewLocalizer(h.bundle, string(event.Locale))

	webhook, err := h.repo.GetWebhook(event.GuildID)
	if err == repository.ErrNoResults {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "WebhookNotSet",
						Other: "A webhook is not set",
					},
				}),
				Flags: uint64(discordgo.MessageFlagsEphemeral),
			},
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching webhook: %w", err)
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "WebhookURL",
					Other: "Webhook URL: {{.WebhookURL}}",
				},
				TemplateData: map[string]string{
					"WebhookURL": webhook.URL,
				},
			}),
			Flags: uint64(discordgo.MessageFlagsEphemeral),
		},
	}, nil
}
