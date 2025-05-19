-- name: CreateOrUpdateBalance :exec
WITH balance_exists AS (
    SELECT EXISTS (
        SELECT 1
        FROM balances
        WHERE user_id = @user_id@@UUID
    ) AS exists_flag
)
INSERT INTO balances (
    user_id,
    available,
    withdrawn,
    updated_at
)
SELECT
    @user_id AS user_id,
    CASE
    WHEN @operation::VARCHAR(16) = 'DEPOSIT' THEN @amount
    ELSE (
        CASE
            WHEN (SELECT exists_flag FROM balance_exists) THEN 0
            ELSE -1
        END
    )
    END AS available,
    CASE
        WHEN @operation = 'WITHDRAW' THEN @amount
        ELSE 0
    END AS withdrawn,
    @updated_at AS updated_at
ON CONFLICT (user_id) DO UPDATE
SET
    available = balances.available + (
    CASE
        WHEN @operation = 'DEPOSIT' THEN @amount
        ELSE -(@amount)
    END
    ),
    withdrawn = balances.withdrawn + (
    CASE
        WHEN @operation = 'WITHDRAW' THEN @amount
        ELSE 0
    END
    ),
    updated_at = @updated_at;
