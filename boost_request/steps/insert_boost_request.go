package steps

import (
	"fmt"

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
		step.repo.DeleteBoostRequest(step.br)
		return nil
	}, nil
}
