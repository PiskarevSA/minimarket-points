package objects

import (
	"errors"

	"github.com/PiskarevSA/minimarket-points/pkg/damm/damm"
)

type OrderNumber string

var NilOrder OrderNumber = ""

var ErrInvalidOrderNumber = errors.New("invalid order number")

func NewOrderNumber(value string) (OrderNumber, error) {
	ok, err := damm.Verify(value)
	if err != nil {
		return NilOrder, ErrInvalidOrderNumber
	}

	if !ok {
		return NilOrder, ErrInvalidOrderNumber
	}

	return OrderNumber(value), nil
}

func (o OrderNumber) String() string {
	return string(o)
}
