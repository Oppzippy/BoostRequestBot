package steps

import (
	"fmt"
	"log"

	"github.com/oppzippy/BoostRequestBot/boost_request/messenger"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type sendPreferredAdvertiserReminderStep struct {
	messenger *messenger.BoostRequestMessenger
	br        *repository.BoostRequest
	repo      repository.Repository
}

func NewSendPreferredAdvertiserReminderStep(repo repository.Repository, messenger *messenger.BoostRequestMessenger, br *repository.BoostRequest) *sendPreferredAdvertiserReminderStep {
	return &sendPreferredAdvertiserReminderStep{
		messenger: messenger,
		br:        br,
		repo:      repo,
	}
}

func (step *sendPreferredAdvertiserReminderStep) Apply() (RevertFunction, error) {
	if len(step.br.PreferredAdvertiserIDs) == 0 {
		return revertNoOp, nil
	}

	delayedMessage, err := step.sendPreferredAdvertiserReminder()
	if err != nil {
		return revertNoOp, err
	}
	return func() error {
		return step.messenger.CancelDelayedMessage(delayedMessage.ID)
	}, nil
}

func (step *sendPreferredAdvertiserReminderStep) sendPreferredAdvertiserReminder() (*repository.DelayedMessage, error) {
	delayedMessage, errChannel := step.messenger.SendPreferredAdvertiserReminder(step.br)
	select {
	case err := <-errChannel:
		return delayedMessage, fmt.Errorf("error sending preferred advertiser reminder: %v", err)
	default:
		err := step.repo.InsertBoostRequestDelayedMessage(step.br, delayedMessage)
		if err != nil {
			return delayedMessage, err
		}
		go func() {
			err, ok := <-errChannel
			if err != nil && ok {
				log.Printf("error sending preferred advertiser reminder: %v", err)
			}
		}()
		return delayedMessage, nil
	}
}
