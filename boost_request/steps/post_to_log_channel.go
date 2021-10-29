package steps

import (
	"log"

	"github.com/oppzippy/BoostRequestBot/boost_request/messenger"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type postToLogChannelStep struct {
	repo      repository.LogChannelRepository
	br        *repository.BoostRequest
	messenger *messenger.BoostRequestMessenger
}

func NewPostToLogChannelStep(repo repository.LogChannelRepository, br *repository.BoostRequest, messenger *messenger.BoostRequestMessenger) *postToLogChannelStep {
	return &postToLogChannelStep{
		repo:      repo,
		br:        br,
		messenger: messenger,
	}
}

func (step *postToLogChannelStep) Apply() (RevertFunction, error) {
	logChannel, err := step.repo.GetLogChannel(step.br.Channel.GuildID)
	if err != repository.ErrNoResults {
		if err != nil {
			log.Printf("Error fetching log channel: %v", err)
		} else {
			_, err := step.messenger.SendLogChannelMessage(step.br, logChannel)
			if err != nil {
				log.Printf("Error sending log channel message: %v", err)
			}
		}
	}
	// We don't want to revert everything else if this fails so never return an error
	return revertNoOp, nil
}
