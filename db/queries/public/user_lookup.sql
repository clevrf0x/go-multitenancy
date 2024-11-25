-- name: CreateUserLookup :one
INSERT INTO user_lookup (tenant_id, email) VALUES ($1, $2) RETURNING *;

-- name: GetUserLookupsByTenant :many
SELECT * FROM user_lookup WHERE tenant_id = $1;

-- name: GetUserLookupByEmail :one
SELECT * FROM user_lookup WHERE email = $1;

-- name: DeleteUserLookupByTenantAndEmail :exec
DELETE FROM user_lookup
WHERE tenant_id = $1 AND email = $2;
