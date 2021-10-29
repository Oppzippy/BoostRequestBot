package steps

type RevertFunction = func() error
type RevertableStep interface {
	Apply() (revert RevertFunction, err error)
}

func revertNoOp() error {
	return nil
}
