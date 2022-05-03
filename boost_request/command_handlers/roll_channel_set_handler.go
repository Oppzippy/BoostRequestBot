package command_handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type RollChannelSetHandler struct {
	discord *discordgo.Session
	bundle  *i18n.Bundle
	repo    repository.Repository
}

func NewRollChannelSetHandler(bundle *i18n.Bundle, repo repository.Repository, discord *discordgo.Session) *RollChannelSetHandler {
	return &RollChannelSetHandler{
		bundle:  bundle,
		repo:    repo,
		discord: discord,
	}
}

func (h *RollChannelSetHandler) Handle(event *discordgo.InteractionCreate, options map[string]*discordgo.ApplicationCommandInteractionDataOption) (*discordgo.InteractionResponse, error) {
	localizer := i18n.NewLocalizer(h.bundle, string(event.Locale))

	channel := options["channel"].ChannelValue(h.discord)
	if channel.Type != discordgo.ChannelTypeGuildText {
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

	err := h.repo.InsertRollChannel(event.GuildID, channel.ID)
	if err != nil {
		return nil, fmt.Errorf("error setting boost request roll channel: %w", err)
	}
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "BoostRequestRollChannelSet",
					Other: "Boost request roll channel set to {.Channel}.",
				},
				TemplateData: map[string]string{
					"Channel": channel.Mention(),
				},
			}),
			Flags: uint64(discordgo.MessageFlagsEphemeral),
		},
	}, nil
}
