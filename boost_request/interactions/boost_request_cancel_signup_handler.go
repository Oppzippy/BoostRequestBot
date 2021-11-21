package interactions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/oppzippy/BoostRequestBot/boost_request/boost_request_manager"
)

type BoostRequestCancelSignUpHandler struct {
	brm *boost_request_manager.BoostRequestManager
}

func NewBoostRequestCancelSignUpHandler(brm *boost_request_manager.BoostRequestManager) *BoostRequestCancelSignUpHandler {
	return &BoostRequestCancelSignUpHandler{
		brm: brm,
	}
}

func (h *BoostRequestCancelSignUpHandler) Matches(discord *discordgo.Session, event *discordgo.InteractionCreate) bool {
	return event.Type == discordgo.InteractionMessageComponent &&
		event.MessageComponentData().CustomID == "boostRequest:cancelSignUp" &&
		event.Member != nil &&
		event.Member.User != nil
}

func (h *BoostRequestCancelSignUpHandler) Handle(discord *discordgo.Session, event *discordgo.InteractionCreate) error {
	removed := h.brm.RemoveAdvertiserFromBoostRequest(event.Message.ID, event.Member.User.ID)
	var content string
	if removed {
		content = "Your signup has been canceled."
	} else {
		content = "You are not signed up for this boost request."
	}
	err := discord.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   1 << 6, // Ephemeral
		},
	})
	return err
}
