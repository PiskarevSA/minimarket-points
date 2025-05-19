package objects

import (
	"github.com/PiskarevSA/minimarket-points/pkg/damm/damm"
)

type OrderNumber string

type OrderNumberError Error

func (e OrderNumberError) Error() string {
	return e.Message
}

var NilOrderNumber OrderNumber = ""

var (
	ErrInvalidOrderNumber    = &OrderNumberError{"invalid order number"}
	ErrInvalidOrderNumberLen = &OrderNumberError{"invalid order number len"}
)

const orderNubmerLen = 10

func NewOrderNumber(value string) (OrderNumber, error) {
	if len(value) != orderNubmerLen {
		return NilOrderNumber, ErrInvalidOrderNumberLen
	}

	ok, err := damm.Verify(value)
	if err != nil {
		return NilOrderNumber, ErrInvalidOrderNumber
	}

	if !ok {
		return NilOrderNumber, ErrInvalidOrderNumber
	}

	return OrderNumber(value), nil
}

func (o OrderNumber) String() string {
	return string(o)
}
