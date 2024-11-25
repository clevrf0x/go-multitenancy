-- name: CreateTenant :one
INSERT INTO tenants (id, name, workspace) VALUES ($1, $2, $3)
RETURNING id, name, workspace, created_at, updated_at;

-- name: GetTenantByID :one
SELECT * FROM tenants WHERE id = $1 AND deleted_at IS NULL;

-- name: GetTenantByWorkspace :one
SELECT * FROM tenants WHERE workspace = $1 AND deleted_at IS NULL;

-- name: ListAllTenants :many
SELECT id, name, workspace, created_at, updated_at FROM tenants WHERE deleted_at IS NULL;

-- name: UpdateTenant :one
UPDATE tenants
SET 
    name = COALESCE($2, name),
    description = COALESCE($3, description),
    logo = COALESCE($4, logo),
    updated_at = (NOW() AT TIME ZONE 'UTC')
WHERE 
    id = $1
    AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteTenant :exec
UPDATE tenants
SET 
    deleted_at = (NOW() AT TIME ZONE 'UTC'),
    updated_at = (NOW() AT TIME ZONE 'UTC')
WHERE 
    id = $1
    AND deleted_at IS NULL;
