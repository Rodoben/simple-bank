-- name: CreateEntry :one
-- description: Create a new entry
INSERT INTO entries (
  account_id,
  amount
)VALUES(
  $1, $2
) RETURNING id;


-- name: GetEntry :one
SELECT * FROM entries WHERE id=$1 LIMIT 1;

-- name: ListEntries :many

SELECT * From entries ORDER BY id LIMIT $1 OFFSET $2;