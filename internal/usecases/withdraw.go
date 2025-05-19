package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/PiskarevSA/minimarket-points/internal/domain/objects"
	"github.com/PiskarevSA/minimarket-points/internal/repo"
)

type Withdraw struct {
	storage AdjustBalanceRepo
}

func NewWithdraw(storage AdjustBalanceRepo) *Withdraw {
	return &Withdraw{
		storage: storage,
	}
}

func (u *Withdraw) parseAndValidateInputs(
	rawOrderNumber string,
	rawAmount string,
) (
	orderNumber objects.OrderNumber,
	amount objects.Amount,
	err error,
) {
	orderNumber, err = objects.NewOrderNumber(rawOrderNumber)
	if err != nil {
		return objects.NilOrderNumber,
			objects.NilAmount,
			&ValidationError{
				Code:    "V1107",
				Message: "invalid order number",
			}
	}

	amount, err = objects.NewAmount(rawAmount)
	if err != nil {
		return objects.NilOrderNumber,
			objects.NilAmount,
			&ValidationError{
				Code:    "V1142",
				Message: "invalid limit value",
			}
	}

	return orderNumber, amount, nil
}

func (u *Withdraw) Do(
	ctx context.Context,
	userId uuid.UUID,
	rawOrderNumber string,
	rawAmount string,
) (proccessedAt time.Time, err error) {
	const op = "withdraw"

	orderNumber, amount, err := u.parseAndValidateInputs(
		rawOrderNumber,
		rawAmount,
	)
	if err != nil {
		return time.Time{}, err
	}

	proccessedAt = timeNow()

	err = u.storage.AdjustBalance(
		ctx,
		userId,
		orderNumber,
		objects.OperationWithdraw,
		amount,
		proccessedAt,
	)
	if err != nil {
		if !errors.Is(err, repo.ErrNotEnoughtBalance) {
			log.Error().
				Err(err).
				Str("layer", "storage").
				Str("op", op).
				Msg("failed to write balance change to storage")

			return time.Time{}, err
		}

		return time.Time{}, &BusinessError{
			Code:    "D1215",
			Message: "insufficient balance",
		}
	}

	return proccessedAt, nil
}
