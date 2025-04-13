-- name: CreateJobRun :one
INSERT INTO job_runs (
  job_id,
  status,
  logs,
  started_at,
  ended_at
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetJobRun :one
SELECT * FROM job_runs
WHERE id = $1 LIMIT 1;

-- name: ListJobRun :many
SELECT * FROM job_runs
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateJobRun :one
UPDATE job_runs
SET
  status = $1,
  logs = $2,
  ended_at = $3
WHERE id = $4
RETURNING *;

-- name: DeleteJobRun :exec
DELETE FROM job_runs
WHERE id = $1;
