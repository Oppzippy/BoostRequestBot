package boost_request

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/oppzippy/BoostRequestBot/boost_request/active_request"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type BoostRequestDiscordHandler struct {
	discord *discordgo.Session
	repo    repository.Repository
	brm     *BoostRequestManager
}

func NewBoostRequestDiscordHandler(discord *discordgo.Session, repo repository.Repository, brm *BoostRequestManager) *BoostRequestDiscordHandler {
	brdh := &BoostRequestDiscordHandler{
		brm:     brm,
		repo:    repo,
		discord: discord,
	}

	discord.Identify.Intents |= discordgo.IntentsGuilds
	discord.Identify.Intents |= discordgo.IntentsGuildMessages
	discord.Identify.Intents |= discordgo.IntentsGuildMessageReactions
	discord.Identify.Intents |= discordgo.IntentsDirectMessages

	discord.AddHandler(brdh.onMessageCreate)
	discord.AddHandler(brdh.onMessageReactionAdd)
	discord.AddHandler(brdh.onMessageReactionRemove)

	return brdh
}

func (brdh *BoostRequestDiscordHandler) onMessageCreate(discord *discordgo.Session, event *discordgo.MessageCreate) {
	if event.Author.ID != discord.State.User.ID && event.GuildID != "" {
		brc, err := brdh.repo.GetBoostRequestChannelByFrontendChannelID(event.GuildID, event.ChannelID)
		if err != nil && err != repository.ErrNoResults {
			log.Printf("Error fetching boost request channel: %v", err)
			return
		}
		if brc != nil {
			var backendMessageID string
			if brc.UsesBuyerMessage {
				backendMessageID = event.ID
			} else {
				err := discord.ChannelMessageDelete(event.ChannelID, event.ID)
				if err != nil {
					log.Printf("Error deleting message: %v", err)
				}
			}
			var embedFields []*repository.MessageEmbedField
			if len(event.Embeds) > 0 {
				embedFields = repository.FromDiscordEmbedFields(event.Embeds[0].Fields)
			}
			_, err = brdh.brm.CreateBoostRequest(brc, BoostRequestPartial{
				RequesterID:      event.Author.ID,
				Message:          event.Content,
				EmbedFields:      embedFields,
				BackendMessageID: backendMessageID,
			})
			if err != nil {
				log.Printf("Error creating boost request: %v", err)
				return
			}
		}
	}
}

func (brdh *BoostRequestDiscordHandler) onMessageReactionAdd(discord *discordgo.Session, event *discordgo.MessageReactionAdd) {
	if event.UserID != discord.State.User.ID {
		br, err := brdh.repo.GetBoostRequestByBackendMessageID(event.ChannelID, event.MessageID)
		if err != nil && err != repository.ErrNoResults {
			log.Printf("Error fetching boost request: %v", err)
			return
		}
		if br != nil {
			switch event.Emoji.Name {
			case AcceptEmoji:
				brdh.brm.AddAdvertiserToBoostRequest(br, event.UserID)
			case StealEmoji:
				brdh.brm.StealBoostRequest(br, event.UserID)
			}
		}
	}
}

func (brdh *BoostRequestDiscordHandler) onMessageReactionRemove(discord *discordgo.Session, event *discordgo.MessageReactionRemove) {
	req, ok := brdh.brm.activeRequests.Load(event.MessageID)
	if ok {
		ar, ok := req.(*active_request.ActiveRequest)
		if ok {
			ar.RemoveSignup(event.UserID)
		}
	}
}
