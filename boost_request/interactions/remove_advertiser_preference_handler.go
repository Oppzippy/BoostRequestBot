package interactions

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/oppzippy/BoostRequestBot/boost_request/boost_request_manager"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

var removeAdvertiserPreferenceRegex = regexp.MustCompile(`^removeAdvertiserPreference:([0-9]+):([A-Fa-f0-9\-]+)$`)

type RemoveAdvertiserPreferenceHandler struct {
	brm  *boost_request_manager.BoostRequestManager
	repo repository.Repository
}

func NewRemoveAdvertiserPreferenceHandler(repo repository.Repository, brm *boost_request_manager.BoostRequestManager) *RemoveAdvertiserPreferenceHandler {
	return &RemoveAdvertiserPreferenceHandler{
		brm:  brm,
		repo: repo,
	}
}

func (h *RemoveAdvertiserPreferenceHandler) Matches(event *discordgo.InteractionCreate) bool {
	_, _, err := h.parseBoostRequestId(event)
	return err == nil
}

func (h *RemoveAdvertiserPreferenceHandler) Handle(discord *discordgo.Session, event *discordgo.InteractionCreate) error {
	guildID, boostRequestID, err := h.parseBoostRequestId(event)
	if err != nil {
		return fmt.Errorf("remove advertiser preference: parsing guid: %v", err)
	}
	br, err := h.repo.GetBoostRequestById(guildID, boostRequestID)
	if err != nil {
		return fmt.Errorf("remove advertiser preference: boost request is nil: guild %v, boost request %v: %v", guildID, boostRequestID, err)
	}

	err = h.brm.CancelBoostRequest(br)
	if err != nil {
		return fmt.Errorf("failed to cancel boost request: %v", err)
	}

	err = discord.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    "**This boost request had a preferred advertiser and was cancelled.**\n\n" + event.Message.Content,
			Embeds:     event.Message.Embeds,
			Components: []discordgo.MessageComponent{},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to respond to interaction: %v", err)
	}

	_, err = h.brm.CreateBoostRequest(&br.Channel, boost_request_manager.BoostRequestPartial{
		RequesterID:   br.RequesterID,
		Message:       br.Message,
		EmbedFields:   br.EmbedFields,
		Price:         br.Price,
		AdvertiserCut: br.AdvertiserCut,
	})
	if err != nil {
		// todo tell them to recreate the br later
		return err
	}
	_, err = discord.FollowupMessageCreate(os.Getenv("DISCORD_APPLICATION_ID"), event.Interaction, true, &discordgo.WebhookParams{
		Content: "A new boost request was created with no advertiser preference.",
	})
	if err != nil {
		return err
	}

	return err
}

func (h *RemoveAdvertiserPreferenceHandler) parseBoostRequestId(event *discordgo.InteractionCreate) (guildID string, boostRequestID uuid.UUID, err error) {
	customID := event.MessageComponentData().CustomID
	matches := removeAdvertiserPreferenceRegex.FindStringSubmatch(customID)
	if matches == nil {
		return "", uuid.UUID{}, errors.New("regex matches is nil")
	}
	boostRequestID, err = uuid.Parse(matches[2])
	return matches[1], boostRequestID, err
}