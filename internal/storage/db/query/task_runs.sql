-- job_run_id UUID REFERENCES job_runs (id) ON DELETE CASCADE,
--     task_id UUID REFERENCES tasks (id) ON DELETE CASCADE,
--     status TEXT CHECK (
--         status IN ('pending', 'running', 'failed', 'completed')
--     ),
--     started_at TIMESTAMP,
--     finished_at TIMESTAMP,
--     exit_code INTEGER,
--     output TEXT,
--     error_message TEXT,

-- name: CreateTaskRun :one
INSERT INTO task_runs (id, job_run_id, task_id, status, started_at, finished_at, exit_code, output, error_message)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, job_run_id, task_id, status, started_at, finished_at, exit_code, output, error_message, created_at;

-- name: GetTaskRun :one
SELECT id, job_run_id, task_id, status, started_at, finished_at, exit_code, output, error_message, created_at
FROM task_runs
WHERE id = $1 LIMIT 1;

-- name: ListTaskRuns :many
SELECT id, job_run_id, task_id, status, started_at, finished_at, exit_code, output, error_message, created_at
FROM task_runs
WHERE job_run_id = $1
ORDER BY created_at DESC;

-- name: ListTaskRunsByJob :many
SELECT tr.id, tr.job_run_id, tr.task_id, tr.status, tr.started_at, tr.finished_at, tr.exit_code, tr.output, tr.error_message, tr.created_at
FROM task_runs tr
JOIN jobs j ON tr.job_run_id = j.id
WHERE j.id = $1
ORDER BY tr.created_at DESC;

-- name: UpdateTaskRun :exec
UPDATE task_runs
SET status = $2, started_at = $3, finished_at = $4, exit_code = $5, output = $6, error_message = $7, updated_at = CURRENT_TIMESTAMP
WHERE id = $1; 

-- name: UpdateTaskRunStatus :exec
UPDATE task_runs
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: UpdateTaskRunError :exec
UPDATE task_runs
SET error_message = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteTaskRun :exec
DELETE FROM task_runs
WHERE id = $1;
