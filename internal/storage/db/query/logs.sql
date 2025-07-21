-- name: CreateLog :exec
INSERT INTO logs (type, job_run_id, task_run_id, timestamp, level, stream, message, metadata)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: CreateSystemLog :exec
INSERT INTO logs (type, level, stream, message, metadata)
VALUES ('system', $1, $2, $3, $4);

-- name: CreateJobLog :exec
INSERT INTO logs (type, job_run_id, level, stream, message, metadata)
VALUES ('job', $1, $2, $3, $4, $5);

-- name: CreateTaskLog :exec
INSERT INTO logs (type, job_run_id, task_run_id, level, stream, message, metadata)
VALUES ('task', $1, $2, $3, $4, $5, $6);

-- name: GetLog :one
SELECT id, type, job_run_id, task_run_id, timestamp, level, stream, message, metadata, created_at
FROM logs
WHERE id = $1 LIMIT 1;

-- name: ListLogs :many
SELECT id, type, job_run_id, task_run_id, timestamp, level, stream, message, metadata, created_at
FROM logs
ORDER BY timestamp DESC
LIMIT $1 OFFSET $2;

-- name: ListLogsByType :many
SELECT id, type, job_run_id, task_run_id, timestamp, level, stream, message, metadata, created_at
FROM logs
WHERE type = $1
ORDER BY timestamp DESC
LIMIT $2 OFFSET $3;

-- name: ListLogsByJobRun :many
SELECT id, type, job_run_id, task_run_id, timestamp, level, stream, message, metadata, created_at
FROM logs
WHERE job_run_id = $1
ORDER BY timestamp ASC;

-- name: ListLogsByTaskRun :many
SELECT id, type, job_run_id, task_run_id, timestamp, level, stream, message, metadata, created_at
FROM logs
WHERE task_run_id = $1
ORDER BY timestamp ASC;

-- name: ListLogsByJobRunPaginated :many
SELECT id, type, job_run_id, task_run_id, timestamp, level, stream, message, metadata, created_at
FROM logs
WHERE job_run_id = $1
ORDER BY timestamp ASC
LIMIT $2 OFFSET $3;

-- name: ListSystemLogs :many
SELECT id, type, job_run_id, task_run_id, timestamp, level, stream, message, metadata, created_at
FROM logs
WHERE type = 'system'
ORDER BY timestamp DESC
LIMIT $1 OFFSET $2;

-- name: ListLogsByLevel :many
SELECT id, type, job_run_id, task_run_id, timestamp, level, stream, message, metadata, created_at
FROM logs
WHERE level = $1
ORDER BY timestamp DESC
LIMIT $2 OFFSET $3;

-- name: ListLogsByJobRunAndLevel :many
SELECT id, type, job_run_id, task_run_id, timestamp, level, stream, message, metadata, created_at
FROM logs
WHERE job_run_id = $1 AND level = $2
ORDER BY timestamp ASC;

-- name: DeleteLog :exec
DELETE FROM logs
WHERE id = $1;

-- name: DeleteLogsByJobRun :exec
DELETE FROM logs
WHERE job_run_id = $1;

-- name: DeleteLogsByTaskRun :exec
DELETE FROM logs
WHERE task_run_id = $1;

-- name: DeleteLogsByType :exec
DELETE FROM logs
WHERE type = $1;

-- name: CountLogsByJobRun :one
SELECT COUNT(*) FROM logs WHERE job_run_id = $1;

-- name: CountLogsByTaskRun :one
SELECT COUNT(*) FROM logs WHERE task_run_id = $1;

-- name: CountLogsByType :one
SELECT COUNT(*) FROM logs WHERE type = $1; 