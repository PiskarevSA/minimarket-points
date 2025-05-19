package objects

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"

	"github.com/PiskarevSA/minimarket-points/pkg/pgx/convtype"
)

type Amount decimal.Decimal

type AmountError Error

func (e AmountError) Error() string {
	return e.Message
}

var NilAmount = Amount{}

var (
	ErrInvaliAmountValue   = &AmountError{"invalid amount value"}
	ErrInvalidAmountFormat = &AmountError{"invalid amount format"}
	ErrInvaliAmountType    = &AmountError{"invalid amount type"}
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
