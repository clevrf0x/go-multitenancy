-- name: CreateUser :one
INSERT INTO users (name, email, password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 AND deleted_at IS NULL;

-- name: UpdateUser :one
UPDATE users
SET name = COALESCE($2, name),
    email = COALESCE($3, email),
    password = COALESCE($4, password),
    last_active_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- Delete a user (soft delete)
-- name: DeleteUser :exec
UPDATE users
SET deleted_at = NOW()
WHERE id = $1;
