-- name: CreateJobLog :one
INSERT INTO job_logs (
  job_id,
  log_level,
  message
) VALUES (
  $1,
  $2,
  $3
) RETURNING *;

-- name: GetJobLog :one
SELECT * FROM job_logs
WHERE id = $1 LIMIT 1;

-- name: GetLogsByJob :many
SELECT * FROM job_logs
WHERE job_id = $1;

-- name: ListJobLogs :many
SELECT * FROM jobs
ORDER BY id DESC
LIMIT $1
OFFSET $2;
