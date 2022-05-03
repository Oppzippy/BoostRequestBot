package command_handlers

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type CreditsAddHandler struct {
	bundle *i18n.Bundle
	repo   repository.Repository
}

func NewCreditsAddHandler(bundle *i18n.Bundle, repo repository.Repository) *CreditsAddHandler {
	return &CreditsAddHandler{
		bundle: bundle,
		repo:   repo,
	}
}

func (h *CreditsAddHandler) Handle(event *discordgo.InteractionCreate, options map[string]*discordgo.ApplicationCommandInteractionDataOption) (*discordgo.InteractionResponse, error) {
	localizer := i18n.NewLocalizer(h.bundle, string(event.Locale))

	user := options["user"].UserValue(nil)
	creditsToAdd := options["credits"].IntValue()

	err := h.repo.AdjustStealCreditsForUser(event.GuildID, user.ID, repository.OperationAdd, int(creditsToAdd))
	if err != nil {
		return nil, fmt.Errorf("error updating steal credits: %w", err)
	}
	newCredits, err := h.repo.GetStealCreditsForUser(event.GuildID, user.ID)
	if err != nil {
		log.Printf("Error fetching steal credits: %v", err)
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "AddedStealCredits",
						One:   "Added {{.Credits}} steal credit.",
						Other: "Added {{.Credits}} steal credits.",
					},
					TemplateData: map[string]int64{
						"Credits": creditsToAdd,
					},
					PluralCount: creditsToAdd,
				}),
				Flags: uint64(discordgo.MessageFlagsEphemeral),
			},
		}, nil
	} else {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "AddedStealCreditsWithNewTotal",
						One:   "Added {{.AddedCredits}} steal credit. New total is {{.TotalCredits}}.",
						Other: "Added {{.AddedCredits}} steal credits. New total is {{.TotalCredits}}.",
					},
					TemplateData: map[string]interface{}{
						"AddedCredits": creditsToAdd,
						"TotalCredits": newCredits,
					},
					PluralCount: creditsToAdd,
				}),
				Flags: uint64(discordgo.MessageFlagsEphemeral),
			},
		}, nil
	}
}
