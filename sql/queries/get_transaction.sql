-- name: GetTransactions :many
SELECT
    id,
    user_id,
    order_number,
    operation,
    amount,
    proccessedAt
FROM transactions
WHERE user_id = @user_id::UUID 
    AND (@operation::VARCHAR(16) = '' OR operation = @operation)
ORDER BY id DESC
OFFSET sqlc.narg('offset')::INTEGER
LIMIT sqlc.narg('limit')::INTEGER;
