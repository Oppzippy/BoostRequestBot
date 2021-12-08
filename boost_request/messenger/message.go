package messenger

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var errDMBlocked = errors.New("the user has the bot blocked")

type messageGenerator interface {
	Message() (*discordgo.MessageSend, error)
}

type message struct {
	dest *MessageDestination
	mg   messageGenerator
}

func newMessage(dest *MessageDestination, mg messageGenerator) *message {
	return &message{
		dest: dest,
		mg:   mg,
	}
}

func (m *message) Send(discord *discordgo.Session) (*discordgo.Message, error) {
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
