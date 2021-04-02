package boost_request

import (
	"database/sql"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

type BoostRequestManager struct {
	discord   *discordgo.Session
	db        *sql.DB
	messenger *BoostRequestMessenger
}

func NewBoostRequestManager(discord *discordgo.Session, db *sql.DB) BoostRequestManager {
	brm := BoostRequestManager{discord: discord, db: db}

	discord.Identify.Intents |= discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions

	discord.AddHandler(brm.onMessageCreate)

	return brm
}

func (brm *BoostRequestManager) onMessageCreate(discord *discordgo.Session, message *discordgo.MessageCreate) {
	brc := brm.GetBoostRequestChannel(message.GuildID, message.ChannelID)
	if brc != nil {
		discord.ChannelMessageDelete(message.ChannelID, message.ID)

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
	err := row.Err()
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		log.Println("Error executing query", err)
		return nil
	}

	brc := BoostRequestChannel{
		GuildID:           guildID,
		FrontendChannelID: frontendChannelID,
	}
	row.Scan(&brc.ID, &brc.BackendChannelID, &brc.UsesBuyerMessage, &brc.NotifiesBuyer)

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
			(boost_request_channel_id, requester_id, backend_message_id, created_at)
			VALUES (?, ?, ?, ?)`,
		brc.ID,
		requesterID,
		message.ID,
		createdAt,
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

	return br, nil
}
