-- name: CreateInventory :one
INSERT INTO inventory (tenant_id, location_id, chicken_count, feed, load_date)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, tenant_id, location_id, chicken_count, feed, load_date, created_at, updated_at;

-- name: GetInventoryByLocationID :one
SELECT id, tenant_id, location_id, chicken_count, feed, load_date, created_at, updated_at
FROM inventory
WHERE tenant_id = $1 AND location_id = $2; -- Get inventory by location ID

-- name: GetInventoryByID :one
SELECT id, tenant_id, location_id, chicken_count, feed, load_date, created_at, updated_at
FROM inventory
WHERE tenant_id = $1 AND id = $2;

-- name: ListInventoryByTenant :many
SELECT id, tenant_id, location_id, chicken_count, feed, load_date, created_at, updated_at
FROM inventory
WHERE tenant_id = $1;

-- name: UpdateInventory :exec
UPDATE inventory
SET chicken_count = $3, feed = $4, load_date = $5, updated_at = now()
WHERE tenant_id = $1 AND id = $2;

-- name: DeleteInventory :exec
DELETE FROM inventory
WHERE tenant_id = $1 AND id = $2;

-- name: GetTotalChickensByTenant :many
SELECT tenant_id, SUM(chicken_count) AS total_chickens
FROM inventory
GROUP BY tenant_id;