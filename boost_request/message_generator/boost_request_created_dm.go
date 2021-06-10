package message_generator

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type BoostRequestCreatedDM struct {
	localizer         *i18n.Localizer
	boostRequest      *repository.BoostRequest
	dmChannelProvider dmChannelProvider
}

func NewBoostRequestCreatedDM(localizer *i18n.Localizer, channelProvider dmChannelProvider, br *repository.BoostRequest) *BoostRequestCreatedDM {
	return &BoostRequestCreatedDM{
		localizer:         localizer,
		boostRequest:      br,
		dmChannelProvider: channelProvider,
	}
}

func (m *BoostRequestCreatedDM) ChannelID() (string, error) {
	channelID, err := m.dmChannelProvider.DMChannel(m.boostRequest.RequesterID)
	if err != nil {
		return "", fmt.Errorf("creating dm channel for boost request created dm: %v", err)
	}
	return channelID, nil
}

func (m *BoostRequestCreatedDM) Message() (*discordgo.MessageSend, error) {
	return &discordgo.MessageSend{}, nil
}
