package command_handlers

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type ChannelsListHandler struct {
	bundle *i18n.Bundle
	repo   repository.Repository
}

func NewChannelsListHandler(bundle *i18n.Bundle, repo repository.Repository) *ChannelsListHandler {
	return &ChannelsListHandler{
		bundle: bundle,
		repo:   repo,
	}
}

func (h *ChannelsListHandler) Handle(event *discordgo.InteractionCreate, options map[string]*discordgo.ApplicationCommandInteractionDataOption) (*discordgo.InteractionResponse, error) {
	localizer := i18n.NewLocalizer(h.bundle, string(event.Locale))

	channels, err := h.repo.GetBoostRequestChannels(event.GuildID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch list of boost request channels: %w", err)
	}

	sb := strings.Builder{}
	for i, brc := range channels {
		options := make([]string, 0, 2)
		if brc.SkipsBuyerDM {
			options = append(options, "doesn't dm buyer")
		}
		if brc.UsesBuyerMessage {
			options = append(options, "reacts directly to buyer's message")
		}
		if len(options) == 0 {
			options = append(options, "none")
		}

		// TODO localize
		sb.WriteString("**Channel ")
		sb.WriteString(fmt.Sprintf("%d", i+1))
		sb.WriteString("**\nFrontend Channel: <#")
		sb.WriteString(brc.FrontendChannelID)
		sb.WriteString(">\nBackend Channel: <#")
		sb.WriteString(brc.BackendChannelID)
		sb.WriteString(">\nOptions: ")
		sb.WriteString(strings.Join(options, ", "))
		sb.WriteString("\n")
	}

	if sb.Len() == 0 {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "NoBoostRequestChannelsExist",
						Other: "There are no boost request channels.",
					},
				}),
				Flags: uint64(discordgo.MessageFlagsEphemeral),
			},
		}, nil
	}
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: sb.String(),
			Flags:   uint64(discordgo.MessageFlagsEphemeral),
		},
	}, nil
}
