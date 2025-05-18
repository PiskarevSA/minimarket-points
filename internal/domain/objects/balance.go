package objects

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Balance struct {
	available Amount
	withdrawn Amount
	updatedAt time.Time
}

var NilBalance = Balance{}

func NewBalance[AmountT string | pgtype.Numeric](
	available AmountT,
	withdrawn AmountT,
	updatedAt time.Time,
) (Balance, error) {
	var (
		balance Balance
		err     error
	)

	balance.available, err = NewAmount(available)
	if err != nil {
		return Balance{}, err
	}

	balance.withdrawn, err = NewAmount(withdrawn)
	if err != nil {
		return Balance{}, err
	}

	balance.updatedAt = updatedAt

	return balance, nil
}

func (b Balance) Available() Amount {
	return b.available
}

func (b Balance) Withdrawn() Amount {
	return b.withdrawn
}

func (b Balance) UpdatedAt() time.Time {
	return b.updatedAt
}

func (b *Balance) SetAvailable(available Amount) {
	b.available = available
}

func (b *Balance) SetWithdrawn(withdrawn Amount) {
	b.withdrawn = withdrawn
}

func (b *Balance) SetUpdatedAt(updatedAt time.Time) {
	b.updatedAt = updatedAt
}
