package postgresql

import (
	"context"
	"errors"

	"github.com/PiskarevSA/minimarket-points/internal/domain/entities"
	"github.com/PiskarevSA/minimarket-points/internal/domain/objects"
	"github.com/PiskarevSA/minimarket-points/internal/dto"
	"github.com/PiskarevSA/minimarket-points/internal/gen/sqlc/postgresql"
	"github.com/PiskarevSA/minimarket-points/internal/repo"
	"github.com/PiskarevSA/minimarket-points/pkg/pgx/convtype"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *PostgreSql) GetTransactions(
	ctx context.Context,
	userId uuid.UUID,
	operation objects.Operation,
	offset int32,
	limit int32,
) ([]entities.Transaction, error) {
	txRows, err := r.querier.GetTransactions(
		ctx,
		postgresql.GetTransactionsParams{
			UserId:    userId,
			Operation: operation.String(),
			Offset:    convtype.Int32ToInt4(offset),
			Limit:     convtype.Int32ToInt4(limit),
		},
	)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}

		return nil, repo.ErrNoTransactionsFound
	}

	return dto.GetTransactionRowsTransactions(txRows), nil
}
