package objects

import "errors"

type Operation string

var NilOperation Operation = ""

const (
	OperationDeposit  Operation = "DEPOSIT"
	OperationWithdraw Operation = "WITHDRAW"
)

var ErrInvalidOperation = errors.New("invalid operation")

func NewOperation(value string) (Operation, error) {
	operation := Operation(value)

	if operation != OperationDeposit &&
		operation != OperationWithdraw {
		return "", ErrInvalidOperation
	}

	return operation, nil
}

func (o Operation) String() string {
	return string(o)
}
