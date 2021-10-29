package sequences

import (
	"fmt"

	"github.com/oppzippy/BoostRequestBot/boost_request/steps"
)

func runSequence(sequenceSteps []steps.RevertableStep) (steps.RevertFunction, error) {
	revertFunctions := make([]steps.RevertFunction, 0, len(sequenceSteps))
	for i, step := range sequenceSteps {
		revert, err := step.Apply()
		revertFunctions = append(revertFunctions, revert)
		if err != nil {
			revertErr := runRevertFunctions(revertFunctions)
			if revertErr != nil {
				return nil, fmt.Errorf("error running sequence step %v: %v, error reverting: %v", i, err, revertErr)
			} else {
				return nil, fmt.Errorf("error running sequence step %v, revert successful: %v", i, err)
			}
		}
	}
	return func() error {
		return runRevertFunctions(revertFunctions)
	}, nil
}

// The functions will be run in reverse order
func runRevertFunctions(reverts []steps.RevertFunction) error {
	for i := len(reverts) - 1; i >= 0; i-- {
		err := reverts[i]()
		if err != nil {
			return err
		}
	}
	return nil
}
