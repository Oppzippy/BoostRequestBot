package command_handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type PrivilegesRemoveHandler struct {
	bundle *i18n.Bundle
	repo   repository.Repository
}

func NewPrivilegesRemoveHandler(bundle *i18n.Bundle, repo repository.Repository) *PrivilegesRemoveHandler {
	return &PrivilegesRemoveHandler{
		bundle: bundle,
		repo:   repo,
	}
}

func (h *PrivilegesRemoveHandler) Handle(event *discordgo.InteractionCreate, options map[string]*discordgo.ApplicationCommandInteractionDataOption) (*discordgo.InteractionResponse, error) {
	localizer := i18n.NewLocalizer(h.bundle, string(event.Locale))
	role := options["role"].RoleValue(nil, event.GuildID)

	privileges, err := h.repo.GetAdvertiserPrivilegesForRole(event.GuildID, role.ID)
	if err != nil && err != repository.ErrNoResults {
		return nil, fmt.Errorf("error fetching advertiser privileges for role: %w", err)
	}
	if privileges == nil {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "RoleHasNoPrivileges",
						Other: "This role has no privileges.",
					},
				}),
				Flags: uint64(discordgo.MessageFlagsEphemeral),
			},
		}, nil
	}
	err = h.repo.DeleteAdvertiserPrivileges(privileges)
	if err != nil {
		return nil, fmt.Errorf("error deleting advertiser privileges: %w", err)
	}
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "RemovedPrivilegesFromRole",
					Other: "Removed privileges from {{.Role}}.",
				},
				TemplateData: map[string]string{
					"Role": role.Mention(),
				},
			}),
			Flags: uint64(discordgo.MessageFlagsEphemeral),
		},
	}, nil
}
