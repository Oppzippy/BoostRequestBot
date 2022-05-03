package command_handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type ChannelsRemoveHandler struct {
	bundle *i18n.Bundle
	repo   repository.Repository
}

func NewChannelsRemoveHandler(bundle *i18n.Bundle, repo repository.Repository) *ChannelsRemoveHandler {
	return &ChannelsRemoveHandler{
		bundle: bundle,
		repo:   repo,
	}
}

func (h *ChannelsRemoveHandler) Handle(event *discordgo.InteractionCreate, options map[string]*discordgo.ApplicationCommandInteractionDataOption) (*discordgo.InteractionResponse, error) {
	localizer := i18n.NewLocalizer(h.bundle, string(event.Locale))

	channel := options["frontend-channel"].ChannelValue(nil)

	brc, err := h.repo.GetBoostRequestChannelByFrontendChannelID(event.GuildID, channel.ID)
	if err == repository.ErrNoResults {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "ChannelIsNotBoostRequestFrontend",
						Other: "{{.Channel}} is not a boost request frontend.",
					},
					TemplateData: map[string]string{
						"Channel": channel.Mention(),
					},
				}),
				Flags: uint64(discordgo.MessageFlagsEphemeral),
			},
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching boost request channel: %w", err)
	}
	err = h.repo.DeleteBoostRequestChannel(brc)
	if err != nil {
		return nil, fmt.Errorf("error deleting boost request channel: %w", err)
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "RemovedBoostRequestChannel",
					Other: "Removed boost request channel {{.Channel}}.",
				},
				TemplateData: map[string]string{
					"Channel": channel.Mention(),
				},
			}),
			Flags: uint64(discordgo.MessageFlagsEphemeral),
		},
	}, nil
}
