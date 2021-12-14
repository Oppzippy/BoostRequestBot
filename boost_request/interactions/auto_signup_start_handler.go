package interactions

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/boost_request_manager"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type AutoSignupEnableHandler struct {
	brm  *boost_request_manager.BoostRequestManager
	repo repository.Repository
}

func NewAutoSignupEnableHandler(repo repository.Repository, brm *boost_request_manager.BoostRequestManager) *AutoSignupEnableHandler {
	return &AutoSignupEnableHandler{
		brm:  brm,
		repo: repo,
	}
}

func (h *AutoSignupEnableHandler) Matches(discord *discordgo.Session, event *discordgo.InteractionCreate) bool {
	return event.Type == discordgo.InteractionApplicationCommand &&
		MatchesCommandPath(event.ApplicationCommandData(), "boostrequest", "autosignup", "start")
}

func (h *AutoSignupEnableHandler) Handle(discord *discordgo.Session, event *discordgo.InteractionCreate, localizer *i18n.Localizer) error {
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

	privileges := h.brm.GetBestRolePrivileges(event.GuildID, event.Member.Roles)

	if privileges.AutoSignupDuration == 0 {
		err := discord.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "NotAllowedToUseCommand",
						Other: "You are not allowed to use this command.",
					},
				}),
				Flags: 1 << 6, // Ephemeral
			},
		})
		return err
	}

	duration := 15 * time.Minute
	options := event.ApplicationCommandData().Options[0].Options[0].Options
	if len(options) == 1 {
		duration = time.Duration(options[0].IntValue()) * time.Minute
		if duration < 1*time.Minute {
			duration = 1 * time.Minute
		} else if duration > time.Duration(privileges.AutoSignupDuration)*time.Second {
			duration = time.Duration(privileges.AutoSignupDuration) * time.Second
		}
	}

	err := h.brm.EnableAutoSignup(event.GuildID, event.Member.User.ID, duration)
	if err != nil {
		return err
	}

	err = discord.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
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
			Flags: 1 << 6, // Ephemeral
		},
	})
	return err
}
