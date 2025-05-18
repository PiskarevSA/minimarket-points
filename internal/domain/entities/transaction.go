package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/PiskarevSA/minimarket-points/internal/domain/objects"
)

type Transaction struct {
	id           int32
	userId       uuid.UUID
	orderNumber  objects.OrderNumber
	operation    objects.Operation
	amount       objects.Amount
	proccessedAt time.Time
}

var NilTransaction = Transaction{}

func NewTransaction[AmountT string | pgtype.Numeric](
	id int32,
	userId uuid.UUID,
	orderNumber string,
	operation string,
	amount AmountT,
	proccessedAt time.Time,
) (Transaction, error) {
	var (
		tx  Transaction
		err error
	)

	tx.orderNumber, err = objects.NewOrderNumber(orderNumber)
	if err != nil {
		return NilTransaction, err
	}

	tx.operation, err = objects.NewOperation(operation)
	if err != nil {
		return NilTransaction, err
	}

	tx.amount, err = objects.NewAmount(amount)
	if err != nil {
		return NilTransaction, err
	}

	tx.userId = userId
	tx.proccessedAt = proccessedAt

	return tx, nil
}

func (t Transaction) Id() int32 {
	return t.id
}

func (t Transaction) Operation() objects.Operation {
	return t.operation
}

func (t Transaction) OrderNumber() objects.OrderNumber {
	return t.orderNumber
}

func (t Transaction) UserId() uuid.UUID {
	return t.userId
}

func (t Transaction) Amount() objects.Amount {
	return t.amount
}

func (t Transaction) ProccessedAt() time.Time {
	return t.proccessedAt
}
