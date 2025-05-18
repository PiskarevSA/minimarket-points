-- name: GetTransactions :many
SELECT
    id,
    user_id,
    order_number,
    operation,
    amount,
    timestamp
FROM transactions
WHERE user_id = @user_id::UUID
ORDER BY id DESC
OFFSET sqlc.narg('offset')::INTEGER
LIMIT sqlc.narg('limit')::INTEGER;
