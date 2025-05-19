package postgresql

import (
	"context"
	"errors"

	"github.com/PiskarevSA/minimarket-points/internal/domain/objects"
	"github.com/PiskarevSA/minimarket-points/internal/dto"
	"github.com/PiskarevSA/minimarket-points/internal/repo"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *PostgreSql) GetBalance(
	ctx context.Context,
	userId uuid.UUID,
) (objects.Balance, error) {

	row, err := r.querier.GetBalance(ctx, userId)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return objects.NilBalance, err
		}

		return objects.NilBalance, repo.ErrNoBalanceFound
	}

	return dto.GetBalanceRowToBalance(row), nil
}
