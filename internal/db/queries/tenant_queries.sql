-- name: GetLocations :many
SELECT id, name, description, created_at, updated_at 
FROM locations;

-- name: GetLocationByID :one
SELECT id, name, description, created_at, updated_at 
FROM locations 
WHERE id = $1;

-- name: CreateLocation :one
INSERT INTO locations (name, description) 
VALUES ($1, $2) 
RETURNING id, name, description, created_at, updated_at;

-- name: UpdateLocation :exec
UPDATE locations 
SET name = $1, description = $2, updated_at = NOW() 
WHERE id = $3;

-- name: DeleteLocation :exec
DELETE FROM locations 
WHERE id = $1;

-- name: GetAllInventory :many
SELECT id, location_id, chicken_count, load_date, created_at, updated_at 
FROM inventory;

-- name: GetInventoryByID :one
SELECT id, location_id, chicken_count, load_date, created_at, updated_at 
FROM inventory 
WHERE id = $1;

-- name: CreateInventory :one
INSERT INTO inventory (location_id, chicken_count, load_date) 
VALUES ($1, $2, $3) 
RETURNING id, location_id, chicken_count, load_date, created_at, updated_at;

-- name: UpdateInventory :exec
UPDATE inventory 
SET location_id = $1, chicken_count = $2, load_date = $3, updated_at = NOW() 
WHERE id = $4;

-- name: DeleteInventory :exec
DELETE FROM inventory 
WHERE id = $1;

-- name: GetAllFeedTypes :many
SELECT id, name, cost_per_unit, created_at, updated_at 
FROM feed_types;

-- name: GetFeedTypeByID :one
SELECT id, name, cost_per_unit, created_at, updated_at 
FROM feed_types 
WHERE id = $1;

-- name: CreateFeedType :one
INSERT INTO feed_types (name, cost_per_unit) 
VALUES ($1, $2) 
RETURNING id, name, cost_per_unit, created_at, updated_at;

-- name: UpdateFeedType :exec
UPDATE feed_types 
SET name = $1, cost_per_unit = $2, updated_at = NOW() 
WHERE id = $3;

-- name: DeleteFeedType :exec
DELETE FROM feed_types 
WHERE id = $1;

-- name: GetAllFeedSchedules :many
SELECT id, feed_type_id, times_per_day, amount_per_feeding, created_at, updated_at 
FROM feed_schedules;

-- name: GetFeedScheduleByID :one
SELECT id, feed_type_id, times_per_day, amount_per_feeding, created_at, updated_at 
FROM feed_schedules 
WHERE id = $1;

-- name: CreateFeedSchedule :one
INSERT INTO feed_schedules (feed_type_id, times_per_day, amount_per_feeding) 
VALUES ($1, $2, $3) 
RETURNING id, feed_type_id, times_per_day, amount_per_feeding, created_at, updated_at;

-- name: UpdateFeedSchedule :exec
UPDATE feed_schedules 
SET feed_type_id = $1, times_per_day = $2, amount_per_feeding = $3, updated_at = NOW() 
WHERE id = $4;

-- name: DeleteFeedSchedule :exec
DELETE FROM feed_schedules 
WHERE id = $1;

-- name: GetAllFeedingLogs :many
SELECT id, feed_type_id, chicken_ids, amount_fed, date_time, comments, created_at, updated_at 
FROM feeding_logs;

-- name: GetFeedingLogByID :one
SELECT id, feed_type_id, chicken_ids, amount_fed, date_time, comments, created_at, updated_at 
FROM feeding_logs 
WHERE id = $1;

-- name: CreateFeedingLog :one
INSERT INTO feeding_logs (feed_type_id, chicken_ids, amount_fed, date_time, comments) 
VALUES ($1, $2, $3, $4, $5) 
RETURNING id, feed_type_id, chicken_ids, amount_fed, date_time, comments, created_at, updated_at;

-- name: UpdateFeedingLog :exec
UPDATE feeding_logs 
SET feed_type_id = $1, chicken_ids = $2, amount_fed = $3, date_time = $4, comments = $5, updated_at = NOW() 
WHERE id = $6;

-- name: DeleteFeedingLog :exec
DELETE FROM feeding_logs 
WHERE id = $1;