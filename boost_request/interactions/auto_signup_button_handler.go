package interactions

import (
	"errors"
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/boost_request_manager"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type AutoSignupButtonHandler struct {
	brm  *boost_request_manager.BoostRequestManager
	repo repository.Repository
}

func NewAutoSignupButtonHandler(repo repository.Repository, brm *boost_request_manager.BoostRequestManager) *AutoSignupButtonHandler {
	return &AutoSignupButtonHandler{
		brm:  brm,
		repo: repo,
	}
}

func (h *AutoSignupButtonHandler) Matches(discord *discordgo.Session, event *discordgo.InteractionCreate) bool {
	if event.Type == discordgo.InteractionMessageComponent {
		_, err := h.parseArgs(event.MessageComponentData().CustomID)
		return err == nil
	}
	return false
}

func (h *AutoSignupButtonHandler) Handle(discord *discordgo.Session, event *discordgo.InteractionCreate, localizer *i18n.Localizer) error {
	guildID := event.GuildID
	if guildID == "" {
		guildID, _ = h.parseArgs(event.MessageComponentData().CustomID)
		if guildID == "" {
			return nil
		}
	}

	member := event.Member
	if member == nil {
		var err error
		member, err = discord.GuildMember(guildID, event.User.ID)
		if err != nil {
			return err
		}
	}

	privileges := h.brm.GetBestRolePrivileges(guildID, member.Roles)

	if privileges.AutoSignupDuration == 0 {
		err := discord.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "NotAllowedToUseButton",
						Other: "You are not allowed to use this button.",
					},
				}),
				Flags: 1 << 6, // Ephemeral
			},
		})
		return err
	}

	duration := time.Duration(privileges.AutoSignupDuration) * time.Second

	err := h.brm.EnableAutoSignup(guildID, member.User.ID, duration)
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

var autoSignupButtonRegex = regexp.MustCompile("^autoSignup(:[0-9]+)?$")

func (h *AutoSignupButtonHandler) parseArgs(customID string) (string, error) {
	matches := autoSignupButtonRegex.FindStringSubmatch(customID)
	if matches == nil {
		return "", errors.New("regex matches is nil")
	}
	return matches[1][1:], nil
}
