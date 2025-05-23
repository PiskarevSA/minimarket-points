// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: create_balance.sql

package postgresql

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createOrUpdateBalance = `-- name: CreateOrUpdateBalance :exec
WITH balance_exists AS (
    SELECT EXISTS (
        SELECT 1
        FROM balances
        WHERE user_id = $1@@UUID
    ) AS exists_flag
)
INSERT INTO balances (
    user_id,
    available,
    withdrawn,
    updated_at
)
SELECT
    $1 AS user_id,
    CASE
    WHEN $2::VARCHAR(16) = 'DEPOSIT' THEN $3
    ELSE (
        CASE
            WHEN (SELECT exists_flag FROM balance_exists) THEN 0
            ELSE -1
        END
    )
    END AS available,
    CASE
        WHEN $2 = 'WITHDRAW' THEN $3
        ELSE 0
    END AS withdrawn,
    $4 AS updated_at
ON CONFLICT (user_id) DO UPDATE
SET
    available = balances.available + (
    CASE
        WHEN $2 = 'DEPOSIT' THEN $3
        ELSE -($3)
    END
    ),
    withdrawn = balances.withdrawn + (
    CASE
        WHEN $2 = 'WITHDRAW' THEN $3
        ELSE 0
    END
    ),
    updated_at = $4
`

type CreateOrUpdateBalanceParams struct {
	UserId    uuid.UUID
	Operation string
	Amount    pgtype.Numeric
	UpdatedAt time.Time
}

func (q *Queries) CreateOrUpdateBalance(ctx context.Context, arg CreateOrUpdateBalanceParams) error {
	_, err := q.db.Exec(ctx, createOrUpdateBalance,
		arg.UserId,
		arg.Operation,
		arg.Amount,
		arg.UpdatedAt,
	)
	return err
}
