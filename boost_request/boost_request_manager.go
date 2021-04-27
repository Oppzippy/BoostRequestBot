package boost_request

import (
	"fmt"
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
	isLoaded       bool
	isLoadedLock   *sync.Mutex
}

func NewBoostRequestManager(discord *discordgo.Session, repo repository.Repository) *BoostRequestManager {
	brm := BoostRequestManager{
		discord:        discord,
		repo:           repo,
		messenger:      NewBoostRequestMessenger(),
		activeRequests: new(sync.Map),
		isLoadedLock:   new(sync.Mutex),
	}

	discord.Identify.Intents |= discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions

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
		if err != nil && err != repository.ErrNoResults {
			log.Printf("Error fetching boost request channel: %v", err)
			return
		}
		if brc != nil {
			if !brc.UsesBuyerMessage {
				err := discord.ChannelMessageDelete(event.ChannelID, event.ID)
				if err != nil {
					log.Printf("Error deleting message: %v", err)
				}
			}
			_, err = brm.CreateBoostRequest(brc, event.Author.ID, event.Message)
			if err != nil {
				log.Printf("Error creating boost request: %v", err)
				return
			}
		}
	}
}

func (brm *BoostRequestManager) onMessageReactionAdd(discord *discordgo.Session, event *discordgo.MessageReactionAdd) {
	if event.UserID != discord.State.User.ID {
		br, err := brm.repo.GetBoostRequestByBackendMessageID(event.ChannelID, event.MessageID)
		if err != nil && err != repository.ErrNoResults {
			log.Printf("Error fetching boost request: %v", err)
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

func (brm *BoostRequestManager) LoadBoostRequests() {
	brm.isLoadedLock.Lock()
	defer brm.isLoadedLock.Unlock()
	if !brm.isLoaded {
		boostRequests, err := brm.repo.GetUnresolvedBoostRequests()
		if err != nil {
			log.Printf("Unable to load boost requests: %v", err)
			return
		}
		brm.isLoaded = true
		for _, br := range boostRequests {
			brm.activeRequests.Store(br.BackendMessageID, NewActiveRequest(*br, brm.setWinner))
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
			return nil, fmt.Errorf("sending backend signup message: %w", err)
		}
	}

	brm.discord.MessageReactionAdd(backendMessage.ChannelID, backendMessage.ID, AcceptEmoji)
	brm.discord.MessageReactionAdd(backendMessage.ChannelID, backendMessage.ID, StealEmoji)

	br.BackendMessageID = backendMessage.ID

	err := brm.repo.InsertBoostRequest(br)
	if err != nil {
		return nil, fmt.Errorf("inserting new boost request in db: %w", err)
	}

	if !brc.SkipsBuyerDM {
		brm.messenger.SendBoostRequestCreatedDM(brm.discord, br)
	}
	brm.activeRequests.Store(br.BackendMessageID, NewActiveRequest(*br, brm.setWinner))

	logChannel, err := brm.repo.GetLogChannel(brc.GuildID)
	if err != repository.ErrNoResults {
		if err != nil {
			log.Printf("Error fetching log channel: %v", err)
		} else {
			brm.messenger.SendLogChannelMessage(brm.discord, br, logChannel)
		}
	}

	return br, nil
}

// Best is defined as the role with the highest weight
func (brm *BoostRequestManager) GetBestRolePrivileges(guildID string, roles []string) *repository.AdvertiserPrivileges {
	var bestPrivileges *repository.AdvertiserPrivileges = nil
	for _, role := range roles {
		privileges, err := brm.repo.GetAdvertiserPrivilegesForRole(guildID, role)
		if err != nil && err != repository.ErrNoResults {
			log.Printf("Error fetching privileges: %v", err)
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
		log.Printf("Error fetching guild member: %v", err)
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
		log.Printf("Error fetching steal credits: %v", err)
		return false
	}
	if credits <= 0 {
		return false
	}

	guildMember, err := brm.discord.GuildMember(br.Channel.GuildID, userID)
	if err != nil {
		log.Printf("Error fetching guild member: %v", err)
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
			log.Printf("Error subtracting boost request credits after use: %v", err)
		}
	}
	return ok
}

func (brm *BoostRequestManager) setWinner(br repository.BoostRequest, userID string) {
	brm.activeRequests.Delete(br.BackendMessageID)
	br.AdvertiserID = userID
	br.IsResolved = true
	br.ResolvedAt = time.Now()
	var rd *repository.RoleDiscount
	// Essentially checking if the user is a bot
	// TODO add IsBot to BoostRequest
	if br.EmbedFields == nil {
		var err error
		rd, err = brm.getRoleDiscountForUser(br.Channel.GuildID, br.RequesterID)
		if err != nil {
			log.Printf("Error searching roles for discount: %v", err)
		}
	}
	br.RoleDiscount = rd
	err := brm.repo.ResolveBoostRequest(&br)
	if err != nil {
		// Log the error but try to keep things running.
		// There will be data loss, but that is better than a lost sale.
		log.Printf("Error resolving boost request: %v", err)
	}

	err = brm.discord.MessageReactionsRemoveAll(br.Channel.BackendChannelID, br.BackendMessageID)
	if err != nil {
		log.Printf("Error removing all reactions: %v", err)
	}
	brm.discord.MessageReactionAdd(br.Channel.BackendChannelID, br.BackendMessageID, ResolvedEmoji)
	_, err = brm.messenger.SendBackendAdvertiserChosenMessage(brm.discord, &br)
	if err != nil {
		log.Printf("Error sending message to boost request backend: %v", err)
	}
	_, err = brm.messenger.SendAdvertiserChosenDMToAdvertiser(brm.discord, &br)
	if err != nil {
		log.Printf("Error sending advertsier chosen DM to advertiser: %v", err)
	}
	if !br.Channel.SkipsBuyerDM {
		_, err = brm.messenger.SendAdvertiserChosenDMToRequester(brm.discord, &br)
		if err != nil {
			log.Printf("Error sending advertiser chosen DM to requester: %v", err)
		}
	}
}

func (brm *BoostRequestManager) getRoleDiscountForUser(guildID, userID string) (*repository.RoleDiscount, error) {
	member, err := brm.discord.GuildMember(guildID, userID)
	if err != nil {
		return nil, err
	}
	var best *repository.RoleDiscount
	for _, role := range member.Roles {
		rd, err := brm.repo.GetRoleDiscountForRole(guildID, role)
		if err == repository.ErrNoResults {
			continue
		}
		if err != nil {
			return nil, err
		}
		if best == nil || rd.Discount.GreaterThan(best.Discount) {
			best = rd
		}
	}
	return best, nil
}

func (brm *BoostRequestManager) signUp(br *repository.BoostRequest, userID string, privileges *repository.AdvertiserPrivileges) {
	req, ok := brm.activeRequests.Load(br.BackendMessageID)
	if !ok {
		return
	}
	r := req.(*activeRequest)
	r.AddSignup(userID, *privileges)
}
