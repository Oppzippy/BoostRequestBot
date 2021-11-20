package interactions

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/oppzippy/BoostRequestBot/boost_request/boost_request_manager"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type BoostRequestSignUpHandler struct {
	brm  *boost_request_manager.BoostRequestManager
	repo repository.Repository
}

func NewBoostRequestSignUpHandler(repo repository.Repository, brm *boost_request_manager.BoostRequestManager) *BoostRequestSignUpHandler {
	return &BoostRequestSignUpHandler{
		brm:  brm,
		repo: repo,
	}
}

func (h *BoostRequestSignUpHandler) Matches(discord *discordgo.Session, event *discordgo.InteractionCreate) bool {
	return event.Type == discordgo.InteractionMessageComponent &&
		event.MessageComponentData().CustomID == "boostRequest:signUp" &&
		event.Member != nil &&
		event.Member.User != nil
}

func (h *BoostRequestSignUpHandler) Handle(discord *discordgo.Session, event *discordgo.InteractionCreate) error {
	br, err := h.repo.GetBoostRequestByBackendMessageID(event.ChannelID, event.Message.ID)
	if err != nil && err != repository.ErrNoResults {
		return fmt.Errorf("error fetching boost request: %v", err)
	}
	if br != nil && !br.IsResolved {
		var content string
		if h.brm.IsAdvertiserSignedUpForBoostRequest(br.BackendMessageID, event.Member.User.ID) {
			content = "You are already signed up for this boost request."
		} else {
			err := h.brm.AddAdvertiserToBoostRequest(br, event.Member.User.ID)

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
