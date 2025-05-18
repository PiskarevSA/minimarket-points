package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/PiskarevSA/minimarket-points/internal/domain/objects"
	"github.com/PiskarevSA/minimarket-points/internal/repo"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type GetBalanceRepo interface {
	GetBalance(
		ctx context.Context,
		userId uuid.UUID,
	) (objects.Balance, error)
}

type GetBalance struct {
	storage GetBalanceRepo
}

func NewGetBalance(storage GetBalanceRepo) *GetBalance {
	return &GetBalance{
		storage: storage,
	}
}

func (u *GetBalance) Do(
	ctx context.Context,
	userId uuid.UUID,
) (balance objects.Balance, err error) {
	const op = "getBalance"

	balance, err = u.storage.GetBalance(ctx, userId)
	if err != nil {
		if !errors.Is(err, repo.ErrNoBalanceFound) {
			log.Error().
				Err(err).
				Str("layer", "usecase").
				Str("op", op).
				Msg("failed to get balance from storage")

			return objects.NilBalance, err
		}

		balance = objects.NilBalance
		balance.SetUpdatedAt(time.Now())

		return balance, err
	}

	return balance, nil
}
