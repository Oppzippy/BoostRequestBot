package messenger

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var errDMBlocked = errors.New("the user has the bot blocked")

type MessageGenerator interface {
	Message() (*discordgo.MessageSend, error)
}

type Message struct {
	dest *MessageDestination
	mg   MessageGenerator
}

func NewMessage(dest *MessageDestination, mg MessageGenerator) *Message {
	return &Message{
		dest: dest,
		mg:   mg,
	}
}

func (m *Message) Send(discord DiscordSender) (*discordgo.Message, error) {
	channelID, err := m.dest.ResolveChannelID(discord)
	if err != nil {
		return nil, fmt.Errorf("resolving channel id: %v", err)
	}
	message, err := m.mg.Message()
	if err != nil {
		return nil, fmt.Errorf("generating message: %v", err)
	}
	sentMessage, err := discord.ChannelMessageSendComplex(channelID, message)

	if err != nil && m.dest.DestinationType == DestinationUser {
		restErr, ok := err.(*discordgo.RESTError)
		if ok && restErr.Message.Code == discordgo.ErrCodeCannotSendMessagesToThisUser {
			return nil, errDMBlocked
		}
	}

	return sentMessage, err
}
