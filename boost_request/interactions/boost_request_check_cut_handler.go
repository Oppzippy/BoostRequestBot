package interactions

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/boost_emojis"
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
		(event.Member != nil || event.User != nil)
}

func (h *BoostRequestCheckCutHandler) Handle(discord *discordgo.Session, event *discordgo.InteractionCreate, localizer *i18n.Localizer) error {
	br, err := h.repo.GetBoostRequestByBackendMessageID(event.Message.ChannelID, event.Message.ID)
	if err != nil {
		return err
	}

	member, err := getGuildMember(discord, event, br.GuildID)
	if err != nil {
		return err
	}

	bestCut := br.AdvertiserCut
	for _, roleID := range member.Roles {
		cut := br.AdvertiserRoleCuts[roleID]
		if cut > bestCut {
			bestCut = cut
		}
	}
	var content string
	if bestCut > 0 {
		emoji := boost_emojis.GoldEmoji
		if event.GuildID != "" {
			emoji, err = h.getGoldEmoji(discord, event.GuildID)
			if err != nil {
				return err
			}
		}
		if br.Discount == 0 {
			content = fmt.Sprintf("Your cut for this boost request is %s.", message_utils.FormatCopperWithEmoji(localizer, bestCut, emoji))
		} else {
			content = fmt.Sprintf(
				"Your discounted cut for this boost request is %s. Before the discount, the cut was %s",
				message_utils.FormatCopperWithEmoji(localizer, bestCut-br.Discount, emoji),
				message_utils.FormatCopperWithEmoji(localizer, bestCut, emoji),
			)
		}
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

func (h *BoostRequestCheckCutHandler) getGoldEmoji(discord *discordgo.Session, guildID string) (string, error) {
	emojis, err := discord.GuildEmojis(guildID)
	if err != nil {
		return "", err
	}
	for _, e := range emojis {
		if strings.ToLower(e.Name) == "gold" {
			return e.MessageFormat(), nil
		}
	}
	return "gold", nil
}
