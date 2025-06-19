-- name: CreateTask :one
INSERT INTO tasks (job_id, name, type, config, "order")
VALUES ($1, $2, $3, $4, $5)
RETURNING id, job_id, name, type, config, "order", created_at;

-- name: CreateBulkTasks :many
INSERT INTO tasks (id, job_id, name, type, config, "order")
VALUES
    ($1, $2, $3, $4, $5, $6),
    ($7, $8, $9, $10, $11, $12)
RETURNING id, job_id, name, type, config, "order", created_at;

-- name: GetTask :one
SELECT id, job_id, name, type, config, "order", created_at
FROM tasks
WHERE id = $1 LIMIT 1;

-- name: GetTasksByJobID :many
SELECT id, job_id, name, type
FROM tasks
WHERE job_id = $1
ORDER BY "order";

-- name: ListTasks :many
SELECT id, job_id, name, type, config, "order", created_at
FROM tasks
WHERE job_id = $1
ORDER BY "order";

-- name: UpdateTask :exec
UPDATE tasks
SET name = $2, type = $3, config = $4, "order" = $5, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = $1;
