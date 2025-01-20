-- name: CreateTenant :one
INSERT INTO public.tenant (name, license_key)
VALUES ($1, $2)
RETURNING id, name, license_key, created_at, updated_at;

-- name: GetTenant :one
SELECT id, name, license_key, created_at, updated_at
FROM public.tenant
WHERE id = $1;

-- name: GetTenantByLicenseKey :one
SELECT id, name, license_key, created_at, updated_at
FROM public.tenant
WHERE license_key = $1;

-- name: Listtenants :many
SELECT id, name, license_key, created_at, updated_at
FROM public.tenant
ORDER BY name;

-- name: UpdateTenant :exec
UPDATE public.tenant
SET name = $1, updated_at = NOW()
WHERE id = $2;

-- name: DeleteTenant :exec
DELETE FROM public.tenant
WHERE id = $1;