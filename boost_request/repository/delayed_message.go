package repository

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type DestinationType string

var (
	DestinationTypeChannel DestinationType = "CHANNEL"
	DestinationTypeUser    DestinationType = "USER"
)

type DelayedMessage struct {
	ID                int64
	DestinationID     string
	DestinationType   DestinationType
	FallbackChannelID string
	Message           *discordgo.MessageSend
	SendAt            time.Time
}
