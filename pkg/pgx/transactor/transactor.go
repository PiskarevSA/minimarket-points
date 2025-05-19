package transactor

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Driver interface {
	BeginTx(ctx context.Context, otps pgx.TxOptions) (pgx.Tx, error)
}

type transactor struct {
	driver Driver
}

func New(driver Driver) Transactor {
	return &transactor{
		driver: driver,
	}
}
