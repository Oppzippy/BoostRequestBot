package interactions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/boost_request_manager"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type BoostRequestCancelSignupHandler struct {
	brm  *boost_request_manager.BoostRequestManager
	repo repository.Repository
}

func NewBoostRequestCancelSignupHandler(repo repository.Repository, brm *boost_request_manager.BoostRequestManager) *BoostRequestCancelSignupHandler {
	return &BoostRequestCancelSignupHandler{
		brm:  brm,
		repo: repo,
	}
}

func (h *BoostRequestCancelSignupHandler) Matches(discord *discordgo.Session, event *discordgo.InteractionCreate) bool {
	return event.Type == discordgo.InteractionMessageComponent &&
		// TODO remove cancelSignUp in later version
		(event.MessageComponentData().CustomID == "boostRequest:cancelSignUp" || event.MessageComponentData().CustomID == "boostRequest:cancelSignup") &&
		(event.Member != nil || event.User != nil)
}

func (h *BoostRequestCancelSignupHandler) Handle(discord *discordgo.Session, event *discordgo.InteractionCreate, localizer *i18n.Localizer) error {
	user := event.User
	if user == nil {
		user = event.Member.User
	}

	br, err := h.repo.GetBoostRequestByBackendMessageID(event.Message.ChannelID, event.Message.ID)
	if err != nil {
		return err
	}
	removed := h.brm.RemoveAdvertiserFromBoostRequest(br, user.ID)
	var content string
	if removed {
		content = "Your signup has been canceled."
	} else {
		content = "You are not signed up for this boost request."
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
