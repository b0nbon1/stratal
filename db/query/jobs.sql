-- name: CreateJob :one
INSERT INTO jobs (
  name,
  schedule,
  type,
  config,
  status,
  retries,
  max_retries
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7
) RETURNING *;

-- name: GetJob :one
SELECT * FROM jobs
WHERE id = $1 LIMIT 1;

-- name: ListJobs :many
SELECT * FROM jobs
ORDER BY id DESC
LIMIT $1
OFFSET $2;

-- name: ListPendingJobs :many
SELECT * FROM jobs
where status = 'pending' AND 
      (schedule IS NOT NULL OR schedule != '')
ORDER BY created_at DESC;

-- name: UpdateJobStatus :one
UPDATE jobs
SET
  status = $2,
  retries = $3,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteJob :exec
DELETE FROM jobs
WHERE id = $1;
