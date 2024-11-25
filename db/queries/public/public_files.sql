-- name: GetPublicFileByID :one
SELECT id, name, path, created_at, updated_at 
FROM files
WHERE id = $1 
  AND deleted_at IS NULL;

-- name: CreatePublicFile :one
INSERT INTO files (
    name,
    path
) VALUES (
    $1, $2
)
RETURNING id, name, path, created_at, updated_at;

-- name: UpdatePublicFile :one
UPDATE files
SET 
    name = COALESCE($2, name),
    path = COALESCE($3, path),
    updated_at = (NOW() AT TIME ZONE 'UTC')
WHERE 
    id = $1
    AND deleted_at IS NULL
RETURNING id, name, path, created_at, updated_at;

-- name: SoftDeletePublicFile :exec
UPDATE files
SET 
    deleted_at = (NOW() AT TIME ZONE 'UTC'),
    updated_at = (NOW() AT TIME ZONE 'UTC')
WHERE 
    id = $1
    AND deleted_at IS NULL;

