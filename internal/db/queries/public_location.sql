-- name: CreateLocation :one
INSERT INTO location (tenant_id, location, zip_code, latitude, longitude, contact_person, phone)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, tenant_id, location, zip_code, latitude, longitude, contact_person, phone, created_at, updated_at;

-- name: GetLocationByID :one
SELECT id, tenant_id, location, zip_code, latitude, longitude, contact_person, phone, created_at, updated_at
FROM location
WHERE tenant_id = $1 AND id = $2;

-- name: GetLocationByName :one
SELECT id, tenant_id, location, zip_code, latitude, longitude, contact_person, phone, created_at, updated_at
FROM location
WHERE tenant_id = $1 AND location = $2;

-- name: ListlocationsByTenant :many
SELECT id, tenant_id, location, zip_code, latitude, longitude, contact_person, phone, created_at, updated_at
FROM location
WHERE tenant_id = $1
ORDER BY location; -- Order by location name

-- name: UpdateLocation :exec
UPDATE location
SET location = $3, zip_code = $4, latitude = $5, longitude = $6, contact_person = $7, phone = $8, updated_at = now()
WHERE tenant_id = $1 AND id = $2;

-- name: DeleteLocation :exec
DELETE FROM location
WHERE tenant_id = $1 AND id = $2;