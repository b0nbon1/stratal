-- name: CreateJob :one
INSERT INTO jobs (name, description, source, raw_payload)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, name, description, source, created_at;

-- name: GetJob :one
SELECT id, user_id, name, description, source, created_at FROM jobs
WHERE id = $1 LIMIT 1;

-- name: UpdateJob :exec
UPDATE jobs
SET name = $2, description = $3, source = $4, raw_payload = $5, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: GetJobWithTasks :one
SELECT j.id, j.user_id, j.name, j.description, j.source, j.created_at,
       json_agg(t.*) AS tasks
FROM jobs j
LEFT JOIN tasks t ON j.id = t.job_id
WHERE j.id = $1
GROUP BY j.id;


-- name: ListJobs :many
SELECT id, user_id, name, description, source, created_at
FROM jobs
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: DeleteJob :exec
DELETE FROM jobs
WHERE id = $1;
