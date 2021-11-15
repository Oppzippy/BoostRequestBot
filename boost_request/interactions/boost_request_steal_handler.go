package interactions

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/oppzippy/BoostRequestBot/boost_request/boost_request_manager"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type BoostRequestStealHandler struct {
	brm  *boost_request_manager.BoostRequestManager
	repo repository.Repository
}

func NewBoostRequestStealHandler(repo repository.Repository, brm *boost_request_manager.BoostRequestManager) *BoostRequestStealHandler {
	return &BoostRequestStealHandler{
		brm:  brm,
		repo: repo,
	}
}

func (h *BoostRequestStealHandler) Matches(discord *discordgo.Session, event *discordgo.InteractionCreate) bool {
	return event.Type == discordgo.InteractionMessageComponent &&
		event.MessageComponentData().CustomID == "boostRequest:steal" &&
		event.Member != nil &&
		event.Member.User != nil
}

func (h *BoostRequestStealHandler) Handle(discord *discordgo.Session, event *discordgo.InteractionCreate) error {
	br, err := h.repo.GetBoostRequestByBackendMessageID(event.ChannelID, event.Message.ID)
	if err != nil && err != repository.ErrNoResults {
		return fmt.Errorf("error fetching boost request: %v", err)
	}
	if br != nil && !br.IsResolved {
		_, usedCredits := h.brm.StealBoostRequest(br, event.Member.User.ID)
		newCredits, err := h.repo.GetStealCreditsForUser(event.GuildID, event.Member.User.ID)
		if err != nil {
			log.Printf("Error fetching steal credits for user: %v", err)
		}
		if usedCredits {
			err = discord.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("You now have %v boost request steal credits.", newCredits),
					Flags:   1 << 6, // Ephemeral
				},
			})
			return err
		} else {
			err = discord.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("No boost request steal credits were spent. You have %v credits.", newCredits),
					Flags:   1 << 6, // Ephemeral
				},
			})
			return err
		}
	}
	return nil
}
