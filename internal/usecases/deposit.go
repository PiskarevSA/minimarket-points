package usecases

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Deposit struct {
}

func NewDeposit() *Deposit {
	return nil
}

func (u *Deposit) Do(
	ctx context.Context,
	userId uuid.UUID,
	rawOrderNumber *string,
	rawAmount string,
) (proccessedAt time.Time, err error) {
	return time.Now(), nil
}
