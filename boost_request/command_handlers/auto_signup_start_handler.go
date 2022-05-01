package command_handlers

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/boost_request_manager"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type AutoSignupEnableHandler struct {
	bundle *i18n.Bundle
	brm    *boost_request_manager.BoostRequestManager
	repo   repository.Repository
}

func NewAutoSignupEnableHandler(bundle *i18n.Bundle, repo repository.Repository, brm *boost_request_manager.BoostRequestManager) *AutoSignupEnableHandler {
	return &AutoSignupEnableHandler{
		bundle: bundle,
		brm:    brm,
		repo:   repo,
	}
}

func (h *AutoSignupEnableHandler) Handle(event *discordgo.InteractionCreate, options map[string]*discordgo.ApplicationCommandInteractionDataOption) (*discordgo.InteractionResponse, error) {
	localizer := i18n.NewLocalizer(h.bundle, string(event.Locale))
	privileges := h.brm.GetBestRolePrivileges(event.GuildID, event.Member.Roles)

	if privileges == nil || privileges.AutoSignupDuration == 0 {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "NotAllowedToUseCommand",
						Other: "You are not allowed to use this command.",
					},
				}),
				Flags: uint64(discordgo.MessageFlagsEphemeral),
			},
		}, nil
	}

	maxDuration := time.Duration(privileges.AutoSignupDuration) * time.Second
	duration := maxDuration
	if durationArgument, ok := options["duration"]; ok {
		duration = time.Duration(durationArgument.IntValue()) * time.Minute
		if duration < 1*time.Minute {
			duration = 1 * time.Minute
		} else if duration > maxDuration {
			duration = maxDuration
		}
	}

	err := h.brm.EnableAutoSignup(event.GuildID, event.Member.User.ID, duration)
	if err != nil {
		return nil, err
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "AutoSignupEnable",
					Other: "You will automatically sign up for all boost requests for the next {{.Duration}}.",
				},
				TemplateData: map[string]interface{}{
					"Duration": duration,
				},
			}),
			Flags: uint64(discordgo.MessageFlagsEphemeral),
		},
	}, nil
}
