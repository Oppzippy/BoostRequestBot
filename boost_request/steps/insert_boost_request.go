package steps

import (
	"fmt"
	"log"

	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type insertBoostRequestStep struct {
	repo repository.BoostRequestRepository
	br   *repository.BoostRequest
}

func NewInsertBoostRequestStep(repo repository.BoostRequestRepository, br *repository.BoostRequest) *insertBoostRequestStep {
	return &insertBoostRequestStep{
		repo: repo,
		br:   br,
	}
}

func (step *insertBoostRequestStep) Apply() (RevertFunction, error) {
	err := step.repo.InsertBoostRequest(step.br)
	if err != nil {
		return revertNoOp, fmt.Errorf("inserting new boost request in db: %w", err)
	}
	return func() error {
		err := step.repo.DeleteBoostRequest(step.br)
		if err != nil {
			log.Printf("error deleting boost request with id %v: %v", step.br.ID, err)
		}
		return nil
	}, nil
}
