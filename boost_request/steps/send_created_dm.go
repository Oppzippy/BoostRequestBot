package steps

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/oppzippy/BoostRequestBot/boost_request/messenger"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type sendCreatedDMStep struct {
	discord   *discordgo.Session
	messenger *messenger.BoostRequestMessenger
	br        *repository.BoostRequest
}

func NewSendCreatedDMStep(discord *discordgo.Session, messenger *messenger.BoostRequestMessenger, br *repository.BoostRequest) *sendCreatedDMStep {
	return &sendCreatedDMStep{
		discord:   discord,
		messenger: messenger,
		br:        br,
	}
}

func (step *sendCreatedDMStep) Apply() (RevertFunction, error) {
	if step.br.Channel != nil && step.br.Channel.SkipsBuyerDM {
		return revertNoOp, nil
	}

	message, err := step.messenger.SendBoostRequestCreatedDM(step.br)
	if err != nil {
		return revertNoOp, fmt.Errorf("error sending boost request created dm: %v", err)
	}
	return func() error {
		return step.discord.ChannelMessageDelete(message.ChannelID, message.ID)
	}, nil
}
