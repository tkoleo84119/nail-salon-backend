-- name: CreateService :one
INSERT INTO services (
    id,
    name,
    price,
    duration_minutes,
    is_addon,
    is_visible,
    note
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING
    id,
    name,
    price,
    duration_minutes,
    is_addon,
    is_visible,
    is_active,
    note;

-- name: GetServiceByID :one
SELECT
    id,
    name,
    price,
    duration_minutes,
    is_addon,
    is_visible,
    is_active,
    note
FROM services
WHERE id = $1;

-- name: GetServiceDetailById :one
SELECT
    id,
    name,
    price,
    duration_minutes,
    is_addon,
    is_visible,
    is_active,
    note
FROM services
WHERE id = $1;

-- name: GetServiceByIds :many
SELECT
    id,
    name,
    price,
    duration_minutes,
    is_addon,
    is_visible,
    is_active,
    note
FROM services
WHERE id = ANY($1::bigint[]);

-- name: CheckServiceNameExists :one
SELECT EXISTS(
    SELECT 1 FROM services
    WHERE name = $1
) AS exists;

-- name: CheckServiceNameExistsExcluding :one
SELECT EXISTS(
    SELECT 1 FROM services
    WHERE name = $1 AND id != $2
) AS exists;
