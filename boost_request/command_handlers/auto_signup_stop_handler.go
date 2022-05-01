package command_handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/boost_request_manager"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type AutoSignupDisableHandler struct {
	bundle *i18n.Bundle
	repo   repository.Repository
	brm    *boost_request_manager.BoostRequestManager
}

func NewAutoSignupDisableHandler(bundle *i18n.Bundle, repo repository.Repository, brm *boost_request_manager.BoostRequestManager) *AutoSignupDisableHandler {
	return &AutoSignupDisableHandler{
		bundle: bundle,
		brm:    brm,
		repo:   repo,
	}
}

func (h *AutoSignupDisableHandler) Handle(event *discordgo.InteractionCreate, options map[string]*discordgo.ApplicationCommandInteractionDataOption) (*discordgo.InteractionResponse, error) {
	localizer := i18n.NewLocalizer(h.bundle, string(event.Locale))

	isEnabled, err := h.repo.IsAutoSignupEnabled(event.GuildID, event.Member.User.ID)
	if err != nil {
		return nil, err
	}
	if !isEnabled {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "AutoSignupNotEnabled",
						Other: "You do not currently have auto sign up active.",
					},
				}),
				Flags: uint64(discordgo.MessageFlagsEphemeral),
			},
		}, nil
	}

	err = h.brm.CancelAutoSignup(event.GuildID, event.Member.User.ID)
	if err != nil {
		return nil, err
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "AutoSignupDisable",
					Other: "You will no longer automatically sign up for boost requests.",
				},
			}),
			Flags: uint64(discordgo.MessageFlagsEphemeral),
		},
	}, nil
}
