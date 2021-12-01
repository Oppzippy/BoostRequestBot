package steps

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/oppzippy/BoostRequestBot/boost_request/messenger"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type sendMessageStep struct {
	discord    *discordgo.Session
	messenger  *messenger.BoostRequestMessenger
	br         *repository.BoostRequest
	channelIDs map[string]struct{}
}

func NewSendMessageStep(discord *discordgo.Session, messenger *messenger.BoostRequestMessenger, br *repository.BoostRequest, channelIDs map[string]struct{}) *sendMessageStep {
	return &sendMessageStep{
		discord:    discord,
		messenger:  messenger,
		br:         br,
		channelIDs: channelIDs,
	}
}

func (step *sendMessageStep) Apply() (RevertFunction, error) {
	if len(step.br.BackendMessages) != 0 {
		// We are probably using the buyer's message as the backend message
		// Anyway, no need to do anything since it's already sent
		return revertNoOp, nil
	}

	reverts := make([]RevertFunction, 0, len(step.channelIDs))
	var err error
	for channelID := range step.channelIDs {
		var revert RevertFunction
		revert, err = step.send(channelID)
		if err != nil {
			break
		}
		reverts = append(reverts, revert)
	}

	return func() error {
		for i := len(reverts) - 1; i >= 0; i-- {
			err := reverts[i]()
			if err != nil {
				return err
			}
			step.br.BackendMessages = step.br.BackendMessages[:i]
		}
		return nil
	}, err
}

func (step *sendMessageStep) send(channelID string) (RevertFunction, error) {
	isDMToPreferredAdvertiser, err := step.isDMToPreferredAdvertiser(channelID)
	if err != nil {
		return nil, err
	}
	backendMessage, err := step.messenger.SendBackendSignupMessage(step.br, channelID, messenger.BackendSignupMessageButtonConfiguration{
		SignUp:       true,
		Steal:        !isDMToPreferredAdvertiser,
		CancelSignup: !isDMToPreferredAdvertiser,
		CheckMyCut:   true,
	})
	if err != nil {
		return revertNoOp, fmt.Errorf("sending backend signup message: %w", err)
	}
	step.br.BackendMessages = append(step.br.BackendMessages, &repository.BoostRequestBackendMessage{
		ChannelID: backendMessage.ChannelID,
		MessageID: backendMessage.ID,
	})
	return func() error {
		err := step.discord.ChannelMessageDelete(backendMessage.ChannelID, backendMessage.ID)
		return err
	}, nil
}

func (step *sendMessageStep) isDMToPreferredAdvertiser(channelID string) (bool, error) {
	channel, err := step.discord.State.Channel(channelID)
	if err == discordgo.ErrNilState || err == discordgo.ErrStateNotFound {
		channel, err = step.discord.Channel(channelID)
	}
	if err != nil {
		return false, err
	}
	if len(channel.Recipients) == 1 {
		_, isPreferredAdvertiser := step.br.PreferredAdvertiserIDs[channel.Recipients[0].ID]
		return isPreferredAdvertiser, nil
	}
	return false, nil
}
