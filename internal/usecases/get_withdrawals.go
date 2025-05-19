package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/PiskarevSA/minimarket-points/internal/domain/entities"
	"github.com/PiskarevSA/minimarket-points/internal/domain/objects"
	"github.com/PiskarevSA/minimarket-points/internal/repo"
)

type GetWithdrawals struct {
	limitMax int32
	storage  GetTransactionsRepo
}

func NewGetWithdrawals(
	limitMax int32,
	storage GetTransactionsRepo,
) *GetWithdrawals {
	return &GetWithdrawals{
		limitMax: limitMax,
		storage:  storage,
	}
}

func (u *GetWithdrawals) validatePagination(
	rawOffset *int32,
	rawLimit *int32,
) (offset int32, limit int32, err error) {
	offset = 0
	limit = 10

	if rawOffset != nil {
		offset = *rawOffset
	}

	if rawLimit != nil {
		limit = *rawLimit
	}

	if limit < 1 || limit > u.limitMax {
		return 0, 0, &ValidationError{
			Code: "V1019",
			Message: fmt.Sprintf(
				"limit value must be < 1 and > %v",
				u.limitMax,
			),
		}
	}

	return offset, limit, nil
}

func (u *GetWithdrawals) Do(
	ctx context.Context,
	userId uuid.UUID,
	rawOffset *int32,
	rawLimit *int32,
) (txs []entities.Transaction, err error) {
	const op = "getWithdrawals"

	offset, limit, err := u.validatePagination(rawLimit, rawLimit)
	if err != nil {
		return nil, err
	}

	txs, err = u.storage.GetTransactions(
		ctx,
		userId,
		objects.OperationWithdraw,
		offset,
		limit,
	)
	if err != nil {
		if !errors.Is(err, repo.ErrNoTransactionsFound) {
			log.Error().
				Err(err).
				Str("layer", "storage").
				Str("op", op).
				Msg("failed to get transactions from storage")

			return nil, err
		}

		return nil, nil
	}

	return txs, nil
}
