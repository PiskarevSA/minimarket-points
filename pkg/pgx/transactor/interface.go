package transactor

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Transactor interface {
	Transact(
		ctx context.Context,
		opts pgx.TxOptions,
		fn func(ctx context.Context) error,
	) error
}
