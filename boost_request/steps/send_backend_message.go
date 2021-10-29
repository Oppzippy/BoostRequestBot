package steps

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/oppzippy/BoostRequestBot/boost_request/messenger"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type sendMessageStep struct {
	discord   *discordgo.Session
	messenger *messenger.BoostRequestMessenger
	br        *repository.BoostRequest
}

func NewSendMessageStep(discord *discordgo.Session, messenger *messenger.BoostRequestMessenger, br *repository.BoostRequest) *sendMessageStep {
	return &sendMessageStep{
		discord:   discord,
		messenger: messenger,
		br:        br,
	}
}

func (step *sendMessageStep) Apply() (RevertFunction, error) {
	if step.br.BackendMessageID != "" {
		// We are probably using the buyer's message as the backend message
		// Anyway, no need to do anything since it's already sent
		return revertNoOp, nil
	}

	backendMessage, err := step.messenger.SendBackendSignupMessage(step.br)
	if err != nil {
		return revertNoOp, fmt.Errorf("sending backend signup message: %w", err)
	}
	step.br.BackendMessageID = backendMessage.ID
	return func() error {
		err := step.discord.ChannelMessageDelete(backendMessage.ChannelID, backendMessage.ID)
		if err != nil {
			return err
		}
		step.br.BackendMessageID = ""
		return nil
	}, nil
}
