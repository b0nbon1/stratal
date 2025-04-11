-- name: CreateJobRun :one
INSERT INTO job_runs (
  owner,
  balance,
  currency
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetJobRun :one
SELECT * FROM job_runs
WHERE id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM job_runs
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListJobRun :many
SELECT * FROM job_runs
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateJobRun :one
UPDATE job_runs
SET balance = $2
WHERE id = $1
RETURNING *;

-- name: DeleteJobRun :exec
DELETE FROM job_runs
WHERE id = $1;
