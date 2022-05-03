package command_handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type ChannelsAddHandler struct {
	discord *discordgo.Session
	bundle  *i18n.Bundle
	repo    repository.Repository
}

func NewChannelsAddHandler(bundle *i18n.Bundle, repo repository.Repository, discord *discordgo.Session) *ChannelsAddHandler {
	return &ChannelsAddHandler{
		discord: discord,
		bundle:  bundle,
		repo:    repo,
	}
}

func (h *ChannelsAddHandler) Handle(event *discordgo.InteractionCreate, options map[string]*discordgo.ApplicationCommandInteractionDataOption) (*discordgo.InteractionResponse, error) {
	localizer := i18n.NewLocalizer(h.bundle, string(event.Locale))

	frontendChannel := options["frontend-channel"].ChannelValue(h.discord)
	backendChannel := options["backend-channel"].ChannelValue(h.discord)

	if frontendChannel.Type != discordgo.ChannelTypeGuildText || backendChannel.Type != discordgo.ChannelTypeGuildText {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "MustBeTextChannel",
						Other: "The specified channel must be a text channel.",
					},
				}),
				Flags: uint64(discordgo.MessageFlagsEphemeral),
			},
		}, nil
	}

	err := h.repo.InsertBoostRequestChannel(&repository.BoostRequestChannel{
		GuildID:           event.GuildID,
		FrontendChannelID: frontendChannel.ID,
		BackendChannelID:  backendChannel.ID,
		UsesBuyerMessage:  frontendChannel.ID == backendChannel.ID,
		SkipsBuyerDM:      frontendChannel.ID == backendChannel.ID,
	})

	if err != nil {
		return nil, fmt.Errorf("error inserting boost request channel: %w", err)
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "AddedBoostRequestChannel",
					Other: "Added boost request frontend {.Frontend} with backend {.Backend}.",
				},
				TemplateData: map[string]string{
					"Frontend": frontendChannel.Mention(),
					"Backend":  backendChannel.Mention(),
				},
			}),
			Flags: uint64(discordgo.MessageFlagsEphemeral),
		},
	}, nil
}
