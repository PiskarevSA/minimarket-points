package usecases

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/PiskarevSA/minimarket-points/internal/domain/objects"
	"github.com/PiskarevSA/minimarket-points/pkg/pgx/transactor"
)

type (
	Transactor interface {
		transactor.Transactor
	}

	AdjustBalanceInRepo interface {
		AdjustBalanceInTx(
			ctx context.Context,
			userId uuid.UUID,
			orderNumber objects.OrderNumber,
			operation objects.Operation,
			amount objects.Amount,
			updatedAt time.Time,
		) error
	}

	CreateTransactionRepo interface {
		CreateTransactionInTx(
			ctx context.Context,
			userId uuid.UUID,
			orderNumber objects.OrderNumber,
			operation objects.Operation,
			amount objects.Amount,
			timestamp time.Time,
		) error
	}
)

var timeNow = func() time.Time {
	return time.Now()
}
