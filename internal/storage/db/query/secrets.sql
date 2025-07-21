-- name: CreateSecret :one
INSERT INTO secrets (user_id, name, encrypted_value)
VALUES ($1, $2, $3)
RETURNING id, user_id, name, created_at;

-- name: GetSecret :one
SELECT id, user_id, name, encrypted_value, created_at
FROM secrets
WHERE id = $1 AND user_id = $2 LIMIT 1;

-- name: GetSecretByName :one
SELECT id, user_id, name, encrypted_value, created_at
FROM secrets
WHERE name = $1 AND user_id = $2 LIMIT 1;

-- name: ListSecrets :many
SELECT id, user_id, name, created_at
FROM secrets
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateSecret :exec
UPDATE secrets
SET encrypted_value = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND user_id = $3;

-- name: DeleteSecret :exec
DELETE FROM secrets
WHERE id = $1 AND user_id = $2;
