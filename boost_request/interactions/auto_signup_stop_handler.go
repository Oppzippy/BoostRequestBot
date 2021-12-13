package interactions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/boost_request_manager"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type AutoSignUpDisableHandler struct {
	brm  *boost_request_manager.BoostRequestManager
	repo repository.Repository
}

func NewAutoSignUpDisableHandler(repo repository.Repository, brm *boost_request_manager.BoostRequestManager) *AutoSignUpDisableHandler {
	return &AutoSignUpDisableHandler{
		brm:  brm,
		repo: repo,
	}
}

func (h *AutoSignUpDisableHandler) Matches(discord *discordgo.Session, event *discordgo.InteractionCreate) bool {
	return event.Type == discordgo.InteractionApplicationCommand &&
		MatchesCommandPath(event.ApplicationCommandData(), "boostrequest", "autosignup", "stop")
}

func (h *AutoSignUpDisableHandler) Handle(discord *discordgo.Session, event *discordgo.InteractionCreate, localizer *i18n.Localizer) error {
	if event.Member == nil {
		err := discord.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "CommandGuildOnly",
						Other: "This command can only be used in guilds.",
					},
				}),
				Flags: 1 << 6, // Ephemeral
			},
		})
		return err
	}

	isEnabled, err := h.repo.IsAutoSignupEnabled(event.GuildID, event.Member.User.ID)
	if err != nil {
		return err
	}
	if !isEnabled {
		discord.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "AutoSignUpNotEnabled",
						Other: "You do not currently have auto sign up active.",
					},
				}),
				Flags: 1 << 6, // Ephemeral
			},
		})
		return nil
	}

	err = h.brm.CancelAutoSignUp(event.GuildID, event.Member.User.ID)
	if err != nil {
		return err
	}

	err = discord.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "AutoSignUpDisable",
					Other: "You will no longer automatically sign up for boost requests.",
				},
			}),
			Flags: 1 << 6, // Ephemeral
		},
	})
	return err
}
