-- name: CreateJob :one
INSERT INTO jobs (id, user_id, name, description, source, raw_payload)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, user_id, name, description, source, created_at;

-- name: GetJob :one
SELECT id, user_id, name, description, source, created_at FROM jobs
WHERE id = $1 LIMIT 1;

-- name: UpdateJob :exec
UPDATE jobs
SET name = $2, description = $3, source = $4, raw_payload = $5, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteJob :exec
DELETE FROM jobs
WHERE id = $1;
