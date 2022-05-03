package command_handlers

import (
	"fmt"
	"net/url"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type WebhookSetHandler struct {
	bundle *i18n.Bundle
	repo   repository.Repository
}

func NewWebhookSetHandler(bundle *i18n.Bundle, repo repository.Repository) *WebhookSetHandler {
	return &WebhookSetHandler{
		bundle: bundle,
		repo:   repo,
	}
}

func (h *WebhookSetHandler) Handle(event *discordgo.InteractionCreate, options map[string]*discordgo.ApplicationCommandInteractionDataOption) (*discordgo.InteractionResponse, error) {
	localizer := i18n.NewLocalizer(h.bundle, string(event.Locale))
	webhookURL := options["webhook-url"].StringValue()

	u, err := url.ParseRequestURI(webhookURL)
	if err != nil || u.Host == "" || u.Scheme == "" {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "InvalidURL",
						Other: "Invalid URL.",
					},
				}),
				Flags: uint64(discordgo.MessageFlagsEphemeral),
			},
		}, nil
	}

	err = h.repo.InsertWebhook(repository.Webhook{
		GuildID: event.GuildID,
		URL:     webhookURL,
	})
	if err != nil {
		return nil, fmt.Errorf("error inserting webhook: %w", err)
	}
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "WebhookSet",
					Other: "Webhook set.",
				},
			}),
			Flags: uint64(discordgo.MessageFlagsEphemeral),
		},
	}, nil
}
