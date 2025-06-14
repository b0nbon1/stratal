-- name: CreateJobRun :one
INSERT INTO job_runs (id, job_id, status, started_at, finished_at, error_message, triggered_by, metadata)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, job_id, status, started_at, finished_at, error_message, triggered_by, metadata, created_at;

-- name: GetJobRun :one
SELECT id, job_id, status, started_at, finished_at, error_message, triggered_by, metadata, created_at
FROM job_runs
WHERE id = $1 LIMIT 1;

-- name: ListJobRuns :many
SELECT id, job_id, status, started_at, finished_at, error_message, triggered_by, metadata, created_at
FROM job_runs
WHERE job_id = $1
ORDER BY created_at DESC;

-- name: UpdateJobRun :exec
UPDATE job_runs
SET status = $2, started_at = $3, finished_at = $4, error_message = $5, triggered_by = $6, metadata = $7, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: UpdateJobRunStatus :exec
UPDATE job_runs
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: UpdateJobRunError :exec
UPDATE job_runs
SET error_message = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteJobRun :exec
DELETE FROM job_runs
WHERE id = $1;
