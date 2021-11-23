package interactions

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages/message_utils"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type BoostRequestCheckCutHandler struct {
	repo repository.Repository
}

func NewBoostRequestCheckCutHandler(repo repository.Repository) *BoostRequestCheckCutHandler {
	return &BoostRequestCheckCutHandler{
		repo: repo,
	}
}

func (h *BoostRequestCheckCutHandler) Matches(discord *discordgo.Session, event *discordgo.InteractionCreate) bool {
	return event.Type == discordgo.InteractionMessageComponent &&
		event.MessageComponentData().CustomID == "boostRequest:checkCut" &&
		event.Member != nil &&
		event.Member.User != nil
}

func (h *BoostRequestCheckCutHandler) Handle(discord *discordgo.Session, event *discordgo.InteractionCreate, localizer *i18n.Localizer) error {
	br, err := h.repo.GetBoostRequestByBackendMessageID(event.Message.ChannelID, event.Message.ID)
	if err != nil {
		return err
	}
	bestCut := br.AdvertiserCut
	for _, roleID := range event.Member.Roles {
		cut := br.AdvertiserRoleCuts[roleID]
		if cut > bestCut {
			bestCut = cut
		}
	}
	var content string
	if bestCut > 0 {
		content = fmt.Sprintf("Your cut for this boost request is %s.", message_utils.FormatCopper(localizer, bestCut))
	} else {
		content = "Your cut for this boost request is unknown."
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
