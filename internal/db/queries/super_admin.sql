-- name: CreateSuperAdmin :one
INSERT INTO super_admin (username, password_hash)
VALUES ($1, $2)
RETURNING id, username, created_at, updated_at;

-- name: GetSuperAdmin :one
SELECT id, username, created_at, updated_at
FROM super_admin
WHERE id = $1;

-- name: GetSuperAdminByUsername :one
SELECT id, username, created_at, updated_at
FROM super_admin
WHERE username = $1;

-- name: UpdateSuperAdmin :exec
UPDATE super_admin
SET username = $1, password_hash = $2, updated_at = NOW()
WHERE id = $3;

-- name: DeleteSuperAdmin :exec
DELETE FROM super_admin
WHERE id = $1;

-- name: ListSuperAdmins :many
SELECT id, username, created_at, updated_at
FROM super_admin
ORDER BY username;