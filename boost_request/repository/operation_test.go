package repository_test

import (
	"testing"

	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func TestOperationFromString(t *testing.T) {
	op, ok := repository.OperationFromString("+")
	if !ok {
		t.Errorf("expected ok but received not ok")
		return
	}
	if op != repository.OperationAdd {
		t.Errorf("Expected %v, received %v", repository.OperationAdd, op)
		return
	}
}

func TestOperationFromStringFailure(t *testing.T) {
	op, ok := repository.OperationFromString("fail")
	if ok {
		t.Errorf("expected ok to be false but it is %v", ok)
		return
	}
	var nilOp repository.Operation
	if op != nilOp {
		t.Errorf("Expected nil operation, received %v", op)
		return
	}
}
