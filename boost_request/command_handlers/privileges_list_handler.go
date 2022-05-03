package command_handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type PrivilegesListHandler struct {
	bundle *i18n.Bundle
	repo   repository.Repository
}

func NewPrivilegesListHandler(bundle *i18n.Bundle, repo repository.Repository) *PrivilegesListHandler {
	return &PrivilegesListHandler{
		bundle: bundle,
		repo:   repo,
	}
}

func (h *PrivilegesListHandler) Handle(event *discordgo.InteractionCreate, options map[string]*discordgo.ApplicationCommandInteractionDataOption) (*discordgo.InteractionResponse, error) {
	localizer := i18n.NewLocalizer(h.bundle, string(event.Locale))

	allPrivileges, err := h.repo.GetAdvertiserPrivilegesForGuild(event.GuildID)
	if err != nil {
		return nil, fmt.Errorf("error listing all privileges: %w", err)
	}

	sb := strings.Builder{}
	for _, p := range allPrivileges {
		// TODO localize
		sb.WriteString("<@&" + p.RoleID + ">")
		sb.WriteString(" Weight: ")
		sb.WriteString(strconv.FormatFloat(p.Weight, 'f', -1, 64))
		sb.WriteString(", Delay: ")
		sb.WriteString(fmt.Sprintf("%ds", p.Delay))
		if p.AutoSignupDuration > 0 {
			sb.WriteString(", Auto Signup Duration: ")
			sb.WriteString(fmt.Sprintf("%dm", p.AutoSignupDuration/60))
		}
		sb.WriteString("\n")
	}
	if sb.Len() == 0 {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "NoPrivilegesAreSet",
						Other: "No privileges are set.",
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
