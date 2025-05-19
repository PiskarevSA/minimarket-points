package postgresql

import (
	"context"
	"time"

	"github.com/PiskarevSA/minimarket-points/internal/domain/objects"
	"github.com/PiskarevSA/minimarket-points/internal/gen/sqlc/postgresql"
	"github.com/PiskarevSA/minimarket-points/pkg/pgx/transactor"
	"github.com/google/uuid"
)

func (r *PostgreSql) CreateTransactionInTx(
	ctx context.Context,
	userId uuid.UUID,
	orderNumber objects.OrderNumber,
	operation objects.Operation,
	amount objects.Amount,
	proccessedAt time.Time,
) error {
	pgxTx := transactor.TxCtxKey.Value(ctx)
	query := r.querier.WithTx(pgxTx)

	return query.CreateTransaction(
		ctx,
		postgresql.CreateTransactionParams{
			UserId:       userId,
			OrderNumber:  orderNumber.String(),
			Operation:    operation.String(),
			Amount:       amount.Numeric(),
			Proccessedat: proccessedAt,
		},
	)
}
