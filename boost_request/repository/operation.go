package repository

type Operation int

const (
	nilOperation = Operation(iota)
	OperationAdd
	OperationSubtract
	OperationMultiply
	OperationDivide
	OperationSet
)

func OperationFromString(op string) (operation Operation, ok bool) {
	switch op {
	case "+":
		operation = OperationAdd
	case "-":
		operation = OperationSubtract
	case "*":
		operation = OperationMultiply
	case "/":
		operation = OperationDivide
	case "=":
		operation = OperationSet
	}
	return operation, false
}
