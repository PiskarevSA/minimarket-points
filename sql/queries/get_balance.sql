-- name: GetBalance :one
SELECT
    available,
    withdrawn,
    updated_at
FROM balances
WHERE user_id = @user_id;
