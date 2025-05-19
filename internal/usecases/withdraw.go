package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"github.com/PiskarevSA/minimarket-points/internal/domain/objects"
	"github.com/PiskarevSA/minimarket-points/internal/repo"
)

type WithdrawRepo interface {
	CreateOrUpdateBalanceInRepo
	CreateTransactionRepo
}

type Withdraw struct {
	storage    WithdrawRepo
	transactor Transactor
}

func NewWithdraw(
	storage WithdrawRepo,
	transactor Transactor,
) *Withdraw {
	return &Withdraw{
		storage:    storage,
		transactor: transactor,
	}
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

	err = u.transact(
		ctx,
		op,
		userId,
		orderNumber,
		amount,
		proccessedAt,
	)
	if err != nil {
		return time.Time{}, err
	}

	return proccessedAt, nil
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

func (u *Withdraw) transact(
	ctx context.Context,
	op string,
	userId uuid.UUID,
	orderNumber objects.OrderNumber,
	amount objects.Amount,
	proccessedAt time.Time,
) error {
	pgxTxOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}

	transactFn := func(ctx context.Context) error {
		err := u.storage.CreateOrUpdateBalanceInTx(
			ctx,
			userId,
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

				return err
			}

			return &BusinessError{
				Code:    "D1215",
				Message: "insufficient balance",
			}
		}

		err = u.storage.CreateTransactionInTx(
			ctx,
			userId,
			orderNumber,
			objects.OperationWithdraw,
			amount,
			proccessedAt,
		)
		if err != nil {
			log.Error().
				Err(err).
				Str("layer", "storage").
				Str("op", op).
				Msg("failed to write transaction to storage")

			return err
		}

		return nil
	}

	return u.transactor.Transact(ctx, pgxTxOpts, transactFn)
}
