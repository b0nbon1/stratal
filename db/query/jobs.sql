-- name: CreateJob :one
INSERT INTO jobs (
  owner,
  balance,
  currency
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetJob :one
SELECT * FROM jobs
WHERE id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM jobs
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListJobs :many
SELECT * FROM jobs
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateJob :one
UPDATE jobs
SET balance = $2
WHERE id = $1
RETURNING *;

-- name: DeleteJob :exec
DELETE FROM jobs
WHERE id = $1;
