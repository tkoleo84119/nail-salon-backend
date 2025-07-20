-- name: CreateService :one
INSERT INTO services (
    id,
    name,
    price,
    duration_minutes,
    is_addon,
    is_visible,
    is_active,
    note
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetServiceByName :one
SELECT * FROM services WHERE name = $1 LIMIT 1;

-- name: GetServiceByID :one
SELECT * FROM services WHERE id = $1 LIMIT 1;

-- name: CheckServiceNameExistsExcluding :one
SELECT EXISTS(
    SELECT 1 FROM services 
    WHERE name = $1 AND id != $2
) AS exists;