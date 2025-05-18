package objects

import (
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"

	"github.com/PiskarevSA/minimarket-points/pkg/pgx/convtype"
)

type Amount decimal.Decimal

var NilAmount = Amount{}

var (
	ErrInvaliAmountValue   = errors.New("invalid amount value")
	ErrInvalidAmountFormat = errors.New("invalid amount format")
	ErrInvaliAmountType    = errors.New("invalid amount type")
)

func NewAmount[T string | pgtype.Numeric](value T) (Amount, error) {
	var dec decimal.Decimal

	var err error

	switch value := any(value).(type) {
	case string:
		dec, err = decimal.NewFromString(value)
		if err != nil {
			return NilAmount, ErrInvalidAmountFormat
		}

	case pgtype.Numeric:
		dec, err = convtype.NumericToDecimal(value)
		if err != nil {
			return NilAmount, err
		}
	default:
		return NilAmount, ErrInvaliAmountType
	}

	if dec.LessThan(decimal.Zero) {
		return NilAmount, ErrInvaliAmountValue
	}

	return Amount(dec), nil
}

func (o Amount) Decimal() decimal.Decimal {
	return decimal.Decimal(o)
}

func (o Amount) Numeric() pgtype.Numeric {
	dec := o.Decimal()

	return convtype.DecimalToNumeric(dec)
}

func (o Amount) String() string {
	return o.Decimal().String()
}

func (o Amount) Set(amount decimal.Decimal) {
	o = Amount(amount)
}
