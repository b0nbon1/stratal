-- name: CreateJobRun :one
INSERT INTO job_runs (job_id, status, triggered_by)
VALUES ($1, $2, $3)
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

-- name: JobRunsWithTasks :one
SELECT jr.id, jr.job_id, jr.status, jr.started_at, jr.finished_at, jr.error_message, jr.triggered_by, jr.metadata, jr.created_at,
       json_agg(tr.*) AS task_runs
FROM job_runs jr
LEFT JOIN task_runs tr ON jr.id = tr.job_run_id
WHERE jr.id = $1
GROUP BY jr.id;



