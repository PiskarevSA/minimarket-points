package handlers

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/PiskarevSA/minimarket-points/internal/domain/entities"
	"github.com/PiskarevSA/minimarket-points/internal/domain/objects"
	"github.com/PiskarevSA/minimarket-points/internal/gen/oapi"
)

type (
	DepositUsecase interface {
		Do(
			ctx context.Context,
			userId uuid.UUID,
			orderNumber *string,
			amount string,
		) (proccessedAt time.Time, err error)
	}

	WithdrawUsecase interface {
		Do(
			ctx context.Context,
			userId uuid.UUID,
			orderNumber string,
			amount string,
		) (proccessedAt time.Time, err error)
	}

	GetBalanceUsecase interface {
		Do(
			ctx context.Context,
			userId uuid.UUID,
		) (balance objects.Balance, err error)
	}

	GetWithdrawalsUsecase interface {
		Do(
			ctx context.Context,
			userId uuid.UUID,
			offset *int32,
			limit *int32,
		) (txs []entities.Transaction, err error)
	}
)

var _ oapi.StrictServerInterface = (*Handlers)(nil)

type Handlers struct {
	depositUsecase        DepositUsecase
	withdrawUsecase       WithdrawUsecase
	getBalanceUsecase     GetBalanceUsecase
	getWithdrawalsUsecase GetWithdrawalsUsecase
}

func New(
	depositUsecase DepositUsecase,
	withdrawUsecase WithdrawUsecase,
	getBalanceUsecase GetBalanceUsecase,
	getWithdrawalsUsecase GetWithdrawalsUsecase,
) *Handlers {
	return &Handlers{
		depositUsecase:        depositUsecase,
		withdrawUsecase:       withdrawUsecase,
		getBalanceUsecase:     getBalanceUsecase,
		getWithdrawalsUsecase: getWithdrawalsUsecase,
	}
}
