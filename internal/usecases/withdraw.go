package usecases

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Withdraw struct {
	storage AdjustBalanceRepo
}

func NewWithdraw(storage AdjustBalanceRepo) *Withdraw {
	return &Withdraw{
		storage: storage,
	}
}

func (u *Withdraw) Do(
	ctx context.Context,
	userId uuid.UUID,
	rawOrderNumber string,
	rawAmount string,
) (proccessedAt time.Time, err error) {
	const op = "withdraw"

	return time.Now(), nil
}
