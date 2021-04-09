package boost_request

import (
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

const AcceptEmoji = "👍"
const StealEmoji = "⏭"
const ResolvedEmoji = "✅"

type BoostRequestManager struct {
	discord        *discordgo.Session
	repo           repository.Repository
	messenger      *BoostRequestMessenger
	activeRequests *sync.Map
}

func NewBoostRequestManager(discord *discordgo.Session, repo repository.Repository) *BoostRequestManager {
	brm := BoostRequestManager{
		discord:        discord,
		repo:           repo,
		messenger:      NewBoostRequestMessenger(),
		activeRequests: new(sync.Map),
	}

	discord.Identify.Intents |= discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions

	discord.AddHandler(brm.onMessageCreate)
	discord.AddHandler(brm.onMessageReactionAdd)

	return &brm
}

func (brm *BoostRequestManager) Destroy() {
	brm.messenger.Destroy(brm.discord)
}

func (brm *BoostRequestManager) onMessageCreate(discord *discordgo.Session, event *discordgo.MessageCreate) {
	if event.Author.ID != discord.State.User.ID && event.GuildID != "" {
		brc, err := brm.repo.GetBoostRequestChannelByFrontendChannelID(event.GuildID, event.ChannelID)
		if err != nil && err != repository.ErrBoostRequestChannelNotFound {
			log.Println("Error fetching boost request channel", err)
			return
		}
		if brc != nil {
			if !brc.UsesBuyerMessage {
				err := discord.ChannelMessageDelete(event.ChannelID, event.ID)
				if err != nil {
					log.Println("Error deleting message", err)
				}
			}
			_, err = brm.CreateBoostRequest(brc, event.Author.ID, event.Message)
			if err != nil {
				log.Println("Error creating boost request", err)
				return
			}
		}
	}
}

func (brm *BoostRequestManager) onMessageReactionAdd(discord *discordgo.Session, event *discordgo.MessageReactionAdd) {
	if event.UserID != discord.State.User.ID {
		br, err := brm.repo.GetBoostRequestByBackendMessageID(event.ChannelID, event.MessageID)
		if err != nil && err != repository.ErrBoostRequestNotFound {
			log.Println("Error fetching boost request", err)
			return
		}
		if br != nil {
			switch event.Emoji.Name {
			case AcceptEmoji:
				brm.addAdvertiserToBoostRequest(br, event.UserID)
			case StealEmoji:
				brm.stealBoostRequest(br, event.UserID)
			}
		}
	}
}

func (brm *BoostRequestManager) CreateBoostRequest(
	brc *repository.BoostRequestChannel, requesterID string, message *discordgo.Message,
) (*repository.BoostRequest, error) {
	createdAt := time.Now().UTC()

	var embedFields []*discordgo.MessageEmbedField
	if len(message.Embeds) == 1 {
		embedFields = message.Embeds[0].Fields
	}
	br := &repository.BoostRequest{
		Channel:     *brc,
		RequesterID: requesterID,
		Message:     message.Content,
		EmbedFields: repository.FromDiscordEmbedFields(embedFields),
		CreatedAt:   createdAt,
	}

	var backendMessage *discordgo.Message
	if brc.UsesBuyerMessage {
		backendMessage = message
	} else {
		var err error
		backendMessage, err = brm.messenger.SendBackendSignupMessage(brm.discord, br)
		if err != nil {
			return nil, err
		}
	}

	brm.discord.MessageReactionAdd(backendMessage.ChannelID, backendMessage.ID, AcceptEmoji)
	brm.discord.MessageReactionAdd(backendMessage.ChannelID, backendMessage.ID, StealEmoji)

	br.BackendMessageID = backendMessage.ID

	err := brm.repo.InsertBoostRequest(br)
	if err != nil {
		return nil, err
	}

	brm.messenger.SendBoostRequestCreatedDM(brm.discord, br)
	brm.activeRequests.Store(br.BackendMessageID, newActiveRequest(*br, brm.setWinner))

	logChannel, err := brm.repo.GetLogChannel(brc.GuildID)
	if err != nil {
		log.Println("Error fetching log channel", err)
	} else if logChannel != "" {
		brm.messenger.SendLogChannelMessage(brm.discord, br, logChannel)
	}

	return br, nil
}

// Best is defined as the role with the highest weight
func (brm *BoostRequestManager) GetBestRolePrivileges(guildID string, roles []string) *repository.AdvertiserPrivileges {
	var bestPrivileges *repository.AdvertiserPrivileges = nil
	for _, role := range roles {
		privileges, err := brm.repo.GetAdvertiserPrivilegesForRole(guildID, role)
		if err != nil {
			log.Println("Error fetching privileges", err)
		}
		if privileges != nil {
			if bestPrivileges == nil || privileges.Weight > bestPrivileges.Weight {
				bestPrivileges = privileges
			}
		}
	}

	return bestPrivileges
}

func (brm *BoostRequestManager) addAdvertiserToBoostRequest(br *repository.BoostRequest, userID string) {
	// TODO cache roles
	guildMember, err := brm.discord.GuildMember(br.Channel.GuildID, userID)
	if err != nil {
		log.Println("Error fetching guild member", err)
		return
	}
	privileges := brm.GetBestRolePrivileges(br.Channel.GuildID, guildMember.Roles)
	if privileges != nil {
		brm.signUp(br, userID, privileges)
	}
}

func (brm *BoostRequestManager) stealBoostRequest(br *repository.BoostRequest, userID string) (ok bool) {
	credits, err := brm.repo.GetStealCreditsForUser(br.Channel.GuildID, userID)
	if err != nil {
		log.Println("Error fetching steal credits", err)
		return false
	}
	if credits <= 0 {
		return false
	}

	guildMember, err := brm.discord.GuildMember(br.Channel.GuildID, userID)
	if err != nil {
		log.Println("Error fetching guild member", err)
		return false
	}
	privileges := brm.GetBestRolePrivileges(br.Channel.GuildID, guildMember.Roles)
	if privileges == nil {
		return
	}

	req, ok := brm.activeRequests.Load(br.BackendMessageID)
	if !ok {
		return false
	}
	r := req.(*activeRequest)
	endTime := br.CreatedAt.Add(time.Duration(privileges.Delay) * time.Second)
	now := time.Now()
	ok = r.SetAdvertiser(userID)
	// Don't subtract steal credits if the 👍 button would have had the same effect
	if ok && now.Before(endTime) {
		err := brm.repo.AdjustStealCreditsForUser(br.Channel.GuildID, userID, repository.OperationSubtract, 1)
		if err != nil {
			log.Println("Error subtracting boost request credits after use", err)
		}
	}
	return ok
}

func (brm *BoostRequestManager) setWinner(br repository.BoostRequest, userID string) {
	brm.activeRequests.Delete(br.BackendMessageID)
	br.AdvertiserID = userID
	br.IsResolved = true
	br.ResolvedAt = time.Now()
	err := brm.repo.ResolveBoostRequest(&br)
	if err != nil {
		log.Println("Error resolving boost request", err)
	}
	err = brm.discord.MessageReactionsRemoveAll(br.Channel.BackendChannelID, br.BackendMessageID)
	if err != nil {
		log.Println("Error removing all reactions", err)
	}
	brm.discord.MessageReactionAdd(br.Channel.BackendChannelID, br.BackendMessageID, ResolvedEmoji)
	_, err = brm.messenger.SendBackendAdvertiserChosenMessage(brm.discord, &br)
	if err != nil {
		log.Println("Error sending message to boost request backend", err)
	}
	if !br.Channel.SkipsBuyerDM {
		_, err = brm.messenger.SendAdvertiserChosenDMToAdvertiser(brm.discord, &br)
		if err != nil {
			log.Println("Error sending advertsier chosen DM to advertiser", err)
		}
		_, err = brm.messenger.SendAdvertiserChosenDMToRequester(brm.discord, &br)
		if err != nil {
			log.Println("Error sending advertiser chosen DM to requester", err)
		}
	}
}

func (brm *BoostRequestManager) signUp(br *repository.BoostRequest, userID string, privileges *repository.AdvertiserPrivileges) {
	req, ok := brm.activeRequests.Load(br.BackendMessageID)
	if !ok {
		return
	}
	r := req.(*activeRequest)
	r.AddSignup(userID, *privileges)
}
