-- name: GetRole :one
SELECT id, name, permissions
FROM public.role
WHERE name = $1;

-- name: ListRoles :many
SELECT id, name, permissions
FROM public.role
ORDER BY name;