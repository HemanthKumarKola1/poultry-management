-- name: CreateUser :one
INSERT INTO public.user (username, password_hash, role, tenant_id)
VALUES ($1, $2, $3, $4)
RETURNING id, username, role, tenant_id, created_at, updated_at;

-- name: GetUser :one
SELECT id, username, role, tenant_id, created_at, updated_at
FROM public.user
WHERE id = $1;

-- name: GetUserByUsername :one
SELECT id, username, role, password_hash, tenant_id, created_at, updated_at
FROM public.user
WHERE username = $1;

-- name: ListusersByTenant :many
SELECT id, username, role, tenant_id, created_at, updated_at
FROM public.user
WHERE tenant_id = $1
ORDER BY username;

-- name: UpdateUser :exec
UPDATE public.user
SET username = $1, role = $2, updated_at = NOW()
WHERE id = $3;

-- name: DeleteUser :exec
DELETE FROM public.user
WHERE id = $1;
