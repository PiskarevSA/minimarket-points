package postgresql

import (
	"context"
	"errors"
	"time"

	"github.com/PiskarevSA/minimarket-points/internal/domain/objects"
	"github.com/PiskarevSA/minimarket-points/internal/gen/sqlc/postgresql"
	"github.com/PiskarevSA/minimarket-points/internal/repo"
	"github.com/PiskarevSA/minimarket-points/pkg/pgcodes"
	"github.com/PiskarevSA/minimarket-points/pkg/pgx/transactor"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *PostgreSql) CreateOrUpdateBalanceInTx(
	ctx context.Context,
	userId uuid.UUID,
	operation objects.Operation,
	amount objects.Amount,
	updatedAt time.Time,
) error {
	pgxTx := transactor.TxCtxKey.Value(ctx)
	query := r.querier.WithTx(pgxTx)

	err := query.CreateOrUpdateBalance(
		ctx,
		postgresql.CreateOrUpdateBalanceParams{
			UserId:    userId,
			Operation: operation.String(),
			Amount:    amount.Numeric(),
			UpdatedAt: updatedAt,
		},
	)
	if err != nil {
		var pgxErr *pgconn.PgError
		if !errors.As(err, &pgxErr) &&
			!pgcodes.IsCheckViolation(pgxErr.Code) {
			return err
		}

		return repo.ErrNotEnoughtBalance
	}

	return nil
}
