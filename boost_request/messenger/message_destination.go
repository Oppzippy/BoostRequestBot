package messenger

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type DestinationType int

const (
	DestinationChannel DestinationType = iota
	DestinationUser
)

type MessageDestination struct {
	DestinationID     string
	DestinationType   DestinationType
	FallbackChannelID string
}

func (dest *MessageDestination) ResolveChannelID(discord *discordgo.Session) (string, error) {
	switch dest.DestinationType {
	case DestinationChannel:
		return dest.DestinationID, nil
	case DestinationUser:
		channel, err := discord.UserChannelCreate(dest.DestinationID)
		if err != nil {
			return "", fmt.Errorf("creating dm channel: %v", err)
		}
		return channel.ID, nil
	}
	return "", fmt.Errorf("invalid destination type: %d", dest.DestinationType)
}
