package command_handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type PrivilegesSetHandler struct {
	bundle *i18n.Bundle
	repo   repository.Repository
}

func NewPrivilegesSetHandler(bundle *i18n.Bundle, repo repository.Repository) *PrivilegesSetHandler {
	return &PrivilegesSetHandler{
		bundle: bundle,
		repo:   repo,
	}
}

func (h *PrivilegesSetHandler) Handle(event *discordgo.InteractionCreate, options map[string]*discordgo.ApplicationCommandInteractionDataOption) (*discordgo.InteractionResponse, error) {
	localizer := i18n.NewLocalizer(h.bundle, string(event.Locale))

	role := options["role"].RoleValue(nil, event.GuildID)
	weight := options["weight"].FloatValue()
	delaySeconds := int(options["delay-in-seconds"].IntValue())
	var autoSignupDurationMinutes int
	if option, ok := options["auto-signup-duration-in-minutes"]; ok {
		autoSignupDurationMinutes = int(option.IntValue())
	}

	err := h.repo.InsertAdvertiserPrivileges(&repository.AdvertiserPrivileges{
		GuildID:            event.GuildID,
		RoleID:             role.ID,
		Weight:             weight,
		Delay:              delaySeconds,
		AutoSignupDuration: autoSignupDurationMinutes * 60,
	})
	if err != nil {
		return nil, fmt.Errorf("error setting privileges: %w", err)
	}
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "PrivilegesSet",
					Other: "Set privileges for {{.Role}}.",
				},
				TemplateData: map[string]string{
					"Role": role.Mention(),
				},
			}),
			Flags: uint64(discordgo.MessageFlagsEphemeral),
		},
	}, nil
}
