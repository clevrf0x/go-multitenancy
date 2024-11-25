-- name: GetTenantFileByID :one
SELECT id, name, path, created_at, updated_at
FROM tenant_files
WHERE id = $1
  AND deleted_at IS NULL;

-- name: CreateTenantFile :one
INSERT INTO tenant_files (
    name,
    path
) VALUES (
    $1, $2
)
RETURNING id, name, path, created_at, updated_at;

-- name: UpdateTenantFile :one
UPDATE tenant_files
SET 
    name = COALESCE($2, name),
    path = COALESCE($3, path),
    updated_at = (NOW() AT TIME ZONE 'UTC')
WHERE 
    id = $1
    AND deleted_at IS NULL
RETURNING id, name, path, created_at, updated_at;

-- name: SoftDeleteTenantFile :exec
UPDATE tenant_files
SET 
    deleted_at = (NOW() AT TIME ZONE 'UTC'),
    updated_at = (NOW() AT TIME ZONE 'UTC')
WHERE 
    id = $1
    AND deleted_at IS NULL;
