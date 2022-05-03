package boost_request

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/application_commands/admin_commands"
	"github.com/oppzippy/BoostRequestBot/boost_request/application_commands/user_commands"
	"github.com/oppzippy/BoostRequestBot/boost_request/boost_emojis"
	"github.com/oppzippy/BoostRequestBot/boost_request/boost_request_manager"
	"github.com/oppzippy/BoostRequestBot/boost_request/interactions"
	"github.com/oppzippy/BoostRequestBot/boost_request/messenger"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
	"github.com/oppzippy/BoostRequestBot/framework/slash_commands"
)

type BoostRequestDiscordHandler struct {
	discord              *discordgo.Session
	repo                 repository.Repository
	brm                  *boost_request_manager.BoostRequestManager
	handlerRemoves       []func()
	interactionRegistry  *InteractionRegistry
	slashCommandRegistry *slash_commands.SlashCommandRegistry
	messenger            *messenger.BoostRequestMessenger
}

func NewBoostRequestDiscordHandler(
	discord *discordgo.Session,
	repo repository.Repository,
	brm *boost_request_manager.BoostRequestManager,
	bundle *i18n.Bundle,
	messenger *messenger.BoostRequestMessenger,
) *BoostRequestDiscordHandler {
	brdh := &BoostRequestDiscordHandler{
		brm:                  brm,
		repo:                 repo,
		discord:              discord,
		handlerRemoves:       make([]func(), 0),
		interactionRegistry:  NewInteractionRegistry(discord, bundle),
		slashCommandRegistry: slash_commands.NewSlashCommandRegistry(),
		messenger:            messenger,
	}

	discord.StateEnabled = true
	discord.State.TrackChannels = true
	discord.State.TrackEmojis = true

	discord.Identify.Intents |= discordgo.IntentsGuilds
	discord.Identify.Intents |= discordgo.IntentsGuildMessages
	discord.Identify.Intents |= discordgo.IntentsGuildMessageReactions
	discord.Identify.Intents |= discordgo.IntentsGuildEmojis
	discord.Identify.Intents |= discordgo.IntentsDirectMessages

	brdh.handlerRemoves = append(
		brdh.handlerRemoves,
		discord.AddHandler(brdh.onMessageCreate),
		discord.AddHandler(brdh.onMessageReactionAdd),
		discord.AddHandler(brdh.onMessageReactionRemove),
		brdh.slashCommandRegistry.AttachToDiscord(discord),
	)

	brdh.interactionRegistry.AddHandler(interactions.NewRemoveAdvertiserPreferenceHandler(repo, brm))
	brdh.interactionRegistry.AddHandler(interactions.NewBoostRequestStealHandler(repo, brm))
	brdh.interactionRegistry.AddHandler(interactions.NewBoostRequestSignupHandler(repo, brm))
	brdh.interactionRegistry.AddHandler(interactions.NewBoostRequestCancelSignupHandler(repo, brm))
	brdh.interactionRegistry.AddHandler(interactions.NewAutoSignupButtonHandler(repo, brm))

	registerSlashCommandHandlers(brdh.slashCommandRegistry, bundle, repo, discord, brm)

	discord.AddHandler(func(discord *discordgo.Session, event *discordgo.Connect) {
		_, err := discord.ApplicationCommandBulkOverwrite(
			discord.State.User.ID,
			"",
			[]*discordgo.ApplicationCommand{
				// User
				user_commands.BoostRequestCommand,

				// Admin
				admin_commands.ChannelsCommand,
				admin_commands.LogChannelCommand,
				admin_commands.PrivilegesCommand,
				admin_commands.RollChannelCommand,
				admin_commands.CreditsCommand,
				admin_commands.WebhookCommand,
			},
		)
		if err != nil {
			log.Printf("Failed to create application command: %v", err)
		}
	})

	return brdh
}

func (brdh *BoostRequestDiscordHandler) Destroy() {
	brdh.interactionRegistry.Destroy()
	for _, remove := range brdh.handlerRemoves {
		remove()
	}
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
			_, err = brdh.brm.CreateBoostRequest(brc, &boost_request_manager.BoostRequestPartial{
				GuildID:          event.GuildID,
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
			case boost_emojis.AcceptEmoji:
				err := brdh.brm.AddAdvertiserToBoostRequest(br, event.UserID)
				if _, isSignupError := err.(*boost_request_manager.BoostRequestSignupError); err != nil && !isSignupError {
					log.Printf("Error adding advertiser to boost request: %v", err)
				}
			case boost_emojis.StealEmoji:
				_, usedCredits := brdh.brm.StealBoostRequest(br, event.UserID)
				newCredits, err := brdh.repo.GetStealCreditsForUser(event.GuildID, event.UserID)
				if err != nil {
					log.Printf("Error fetching steal credits for user: %v", err)
				}
				if usedCredits {
					brdh.messenger.SendCreditsUpdateDM(event.UserID, newCredits)
				}
			}
		}
	}
}

func (brdh *BoostRequestDiscordHandler) onMessageReactionRemove(discord *discordgo.Session, event *discordgo.MessageReactionRemove) {
	br, err := brdh.repo.GetBoostRequestByBackendMessageID(event.ChannelID, event.MessageID)
	if err != nil {
		log.Printf("Error fetching boost request: %v", err)
		return
	}
	brdh.brm.RemoveAdvertiserFromBoostRequest(br, event.UserID)
}
