package boost_request

import (
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

const ACCEPT_EMOJI = "👍"
const STEAL_EMOJI = "⏭"

type BoostRequestManager struct {
	discord    *discordgo.Session
	repo       repository.Repository
	messenger  *BoostRequestMessenger
	privileges *sync.Map
}

func NewBoostRequestManager(discord *discordgo.Session, repo repository.Repository) *BoostRequestManager {
	brm := BoostRequestManager{
		discord:    discord,
		repo:       repo,
		messenger:  NewBoostRequestMessenger(),
		privileges: new(sync.Map),
	}

	discord.Identify.Intents |= discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions

	discord.AddHandler(brm.onMessageCreate)
	discord.AddHandler(brm.onMessageReactionAdd)

	return &brm
}

func (brm *BoostRequestManager) Destroy() {
	brm.messenger.Destroy(brm.discord)
}

func (brm *BoostRequestManager) onMessageCreate(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if !message.Author.Bot && message.GuildID != "" {
		brc, err := brm.repo.GetBoostRequestChannelByFrontendChannelID(message.GuildID, message.ChannelID)
		if err != nil {
			log.Println("Error fetching boost request channel", err)
			return
		}
		if brc != nil {
			_, err := brm.CreateBoostRequest(brc, message.Author.ID, message.Content)
			if err != nil {
				log.Println("Error creating boost request", err)
				return
			}
			discord.ChannelMessageDelete(message.ChannelID, message.ID)
		}
	}
}

func (brm *BoostRequestManager) onMessageReactionAdd(discord *discordgo.Session, event *discordgo.MessageReactionAdd) {
	br, err := brm.repo.GetBoostRequestByBackendMessageID(event.ChannelID, event.MessageID)
	if err != nil {
		log.Println("Error fetching boost request", err)
		return
	}
	if br != nil {
		switch event.Emoji.Name {
		case ACCEPT_EMOJI:
			brm.addAdvertiserToBoostRequest(br, event.UserID)
		case STEAL_EMOJI:
			// TODO not implemented
		}
	}
}

func (brm *BoostRequestManager) CreateBoostRequest(brc *repository.BoostRequestChannel, requesterID string, request string) (*repository.BoostRequest, error) {
	createdAt := time.Now().UTC()
	var br *repository.BoostRequest
	{
		boostRequest := repository.BoostRequest{
			Channel:     brc,
			RequesterID: requesterID,
			Message:     request,
			CreatedAt:   &createdAt,
		}
		br = &boostRequest
	}

	message, err := brm.messenger.SendBackendSignupMessage(brm.discord, br)
	if err != nil {
		return nil, err
	}

	br.BackendMessageID = message.ID

	err = brm.repo.InsertBoostRequest(br)
	if err != nil {
		return nil, err
	}

	brm.messenger.SendBoostRequestCreatedDM(brm.discord, br)

	return br, nil
}

// Best is defined as the role with the highest weight
func (brm *BoostRequestManager) GetBestRolePrivileges(roles []string) *AdvertiserPrivileges {
	var bestPrivileges *AdvertiserPrivileges = nil
	for _, role := range roles {
		value, ok := brm.privileges.Load(role)
		if ok {
			privileges, ok := value.(*AdvertiserPrivileges)
			if ok && privileges != nil {
				if bestPrivileges == nil || privileges.Weight > bestPrivileges.Weight {
					bestPrivileges = privileges
				}
			}
		}
	}
	privileges := *bestPrivileges
	return &privileges
}

func (brm *BoostRequestManager) addAdvertiserToBoostRequest(br *repository.BoostRequest, userID string) {
	guildMember, err := brm.discord.GuildMember(br.Channel.GuildID, userID)
	if err != nil {
		log.Println("Error fetching guild member", err)
		return
	}
	privileges := brm.GetBestRolePrivileges(guildMember.Roles)
	if privileges != nil {
		// TODO add user to waitlist
	}
}
