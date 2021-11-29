package steps

import (
	"sync"

	"github.com/oppzippy/BoostRequestBot/boost_request/active_request"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type storeActiveRequestStep struct {
	activeRequests *sync.Map
	br             *repository.BoostRequest
	setWinner      func(*active_request.AdvertiserChosenEvent)
}

func NewStoreActiveRequestStep(
	activeRequests *sync.Map,
	br *repository.BoostRequest,
	setWinner func(*active_request.AdvertiserChosenEvent),
) *storeActiveRequestStep {
	return &storeActiveRequestStep{
		activeRequests: activeRequests,
		br:             br,
		setWinner:      setWinner,
	}
}

func (step *storeActiveRequestStep) Apply() (RevertFunction, error) {
	step.activeRequests.Store(step.br.ID, active_request.NewActiveRequest(*step.br, step.setWinner))
	return func() error {
		arInterface, ok := step.activeRequests.LoadAndDelete(step.br.ID)
		if ok {
			ar := arInterface.(*active_request.ActiveRequest)
			ar.Destroy()
		}
		return nil
	}, nil
}
