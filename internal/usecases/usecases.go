package usecases

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/PiskarevSA/minimarket-points/internal/domain/entities"
	"github.com/PiskarevSA/minimarket-points/internal/domain/objects"
	"github.com/PiskarevSA/minimarket-points/pkg/pgx/transactor"
)

type (
	Transactor interface {
		transactor.Transactor
	}

	CreateOrUpdateBalanceInRepo interface {
		CreateOrUpdateBalanceInTx(
			ctx context.Context,
			userId uuid.UUID,
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
			proccessedAt time.Time,
		) error
	}

	GetTransactionsRepo interface {
		GetTransactions(
			ctx context.Context,
			userId uuid.UUID,
			operation objects.Operation,
			offset int32,
			limit int32,
		) ([]entities.Transaction, error)
	}
)

var timeNow = func() time.Time {
	return time.Now()
}
