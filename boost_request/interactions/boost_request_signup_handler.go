package interactions

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/boost_request_manager"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type BoostRequestSignupHandler struct {
	brm  *boost_request_manager.BoostRequestManager
	repo repository.Repository
}

func NewBoostRequestSignupHandler(repo repository.Repository, brm *boost_request_manager.BoostRequestManager) *BoostRequestSignupHandler {
	return &BoostRequestSignupHandler{
		brm:  brm,
		repo: repo,
	}
}

func (h *BoostRequestSignupHandler) Matches(discord *discordgo.Session, event *discordgo.InteractionCreate) bool {
	return event.Type == discordgo.InteractionMessageComponent &&
		event.MessageComponentData().CustomID == "boostRequest:signUp" &&
		(event.User != nil || event.Member != nil)
}

func (h *BoostRequestSignupHandler) Handle(discord *discordgo.Session, event *discordgo.InteractionCreate, localizer *i18n.Localizer) error {
	user := event.User
	if user == nil {
		user = event.Member.User
	}

	br, err := h.repo.GetBoostRequestByBackendMessageID(event.ChannelID, event.Message.ID)
	if err != nil && err != repository.ErrNoResults {
		return fmt.Errorf("error fetching boost request: %v", err)
	}
	if br != nil && !br.IsResolved {
		var content string
		if h.brm.IsAdvertiserSignedUpForBoostRequest(br, user.ID) {
			content = "You are already signed up for this boost request."
		} else {
			err := h.brm.AddAdvertiserToBoostRequest(br, user.ID)

			// We could check if it's a BoostRequestSignupError in addition to the switch to ensure
			// if a case is missed, it doesn't error. An error is what we want though, so it is noticed
			// and a new message can be added here.
			switch err {
			case nil:
				content = "You have been signed up."
			case boost_request_manager.ErrNoPrivileges:
				content = "You do not have permission to sign up for boost requests."
			case boost_request_manager.ErrNotPreferredAdvertiser:
				content = "You are not a preferred advertiser for this boost request."
			default:
				return err
			}
		}
		err = discord.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: content,
				Flags:   1 << 6, // Ephemeral
			},
		})
		return err
	}
	return nil
}
