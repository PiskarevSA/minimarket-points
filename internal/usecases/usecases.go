package usecases

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/PiskarevSA/minimarket-points/internal/domain/objects"
)

type AdjustBalanceRepo interface {
	AdjustBalance(
		ctx context.Context,
		userId uuid.UUID,
		orderNumber objects.OrderNumber,
		operation objects.Operation,
		amount objects.Amount,
		updatedAt time.Time,
	) error
}

var timeNow = func() time.Time {
	return time.Now()
}
