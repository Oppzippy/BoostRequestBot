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
			hasPrivileges, err := h.brm.AddAdvertiserToBoostRequest(br, event.Member.User.ID)
			if err != nil {
				return err
			}
			if hasPrivileges {
				content = "You have been signed up."
			} else {
				content = "You do not have permission to sign up for boost requests."
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
