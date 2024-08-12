-- name: CreateTransfer :one
-- description: Create a new entry
INSERT INTO transfers (
  from_account_id,
  to_account_id,
  amount
)VALUES(
  $1, $2, $3
) RETURNING *;


-- name: GetTransfer :one
SELECT * FROM transfers WHERE id=$1 LIMIT 1;

-- name: ListTransfers :many

SELECT * From transfers ORDER BY id LIMIT $1 OFFSET $2;