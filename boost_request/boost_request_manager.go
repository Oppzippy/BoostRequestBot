package boost_request

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

const ACCEPT_EMOJI = "👍"
const STEAL_EMOJI = "⏭"

type BoostRequestManager struct {
	discord    *discordgo.Session
	db         *sql.DB
	messenger  *BoostRequestMessenger
	privileges *sync.Map
}

func NewBoostRequestManager(discord *discordgo.Session, db *sql.DB) *BoostRequestManager {
	brm := BoostRequestManager{
		discord:    discord,
		db:         db,
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
		brc := brm.GetBoostRequestChannel(message.GuildID, message.ChannelID)
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
	br, err := GetBoostRequestByBackendMessageID(brm.db, event.ChannelID, event.MessageID)
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

func (brm *BoostRequestManager) GetBoostRequestChannel(guildID string, frontendChannelID string) *BoostRequestChannel {
	row := brm.db.QueryRow(
		`SELECT id, backend_channel_id, uses_buyer_message, notifies_buyer
			FROM boost_request_channel
			WHERE guild_id = ? AND frontend_channel_id = ?`,
		guildID,
		frontendChannelID,
	)
	brc := BoostRequestChannel{
		GuildID:           guildID,
		FrontendChannelID: frontendChannelID,
	}
	var usesBuyerMessage, notifiesBuyer int
	err := row.Scan(&brc.ID, &brc.BackendChannelID, &usesBuyerMessage, &notifiesBuyer)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		log.Println("Error executing query", err)
		return nil
	}
	brc.UsesBuyerMessage = usesBuyerMessage != 0
	brc.NotifiesBuyer = notifiesBuyer != 0

	return &brc
}

func (brm *BoostRequestManager) CreateBoostRequest(brc *BoostRequestChannel, requesterID string, request string) (*BoostRequest, error) {
	createdAt := time.Now().UTC()
	var br *BoostRequest
	{
		boostRequest := BoostRequest{
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

	res, err := brm.db.Exec(
		`INSERT INTO boost_request
			(boost_request_channel_id, requester_id, backend_message_id, message, created_at)
			VALUES (?, ?, ?, ?, ?)`,
		brc.ID,
		br.RequesterID,
		message.ID,
		br.Message,
		br.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err == nil {
		br.ID = int(id)
		br.BackendMessageID = message.ID
	} else {
		br, err = GetBoostRequestByBackendMessageID(brm.db, message.ChannelID, message.ID)
		if err != nil {
			log.Println("Inserted boost request but failed to retrieve it immediately after!", err)
			return nil, err
		}
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

func (brm *BoostRequestManager) addAdvertiserToBoostRequest(br *BoostRequest, userID string) {
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
