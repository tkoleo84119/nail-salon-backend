-- name: CreateStore :one
INSERT INTO stores (
    id,
    name,
    address,
    phone,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, NOW(), NOW()
) RETURNING
    id,
    name,
    address,
    phone,
    is_active;

-- name: GetAllActiveStoresName :many
SELECT
    id,
    name
FROM stores
WHERE is_active = true
ORDER BY name;

-- name: GetStoreByID :one
SELECT
    id,
    name,
    is_active
FROM stores
WHERE id = $1;

-- name: GetStoresByIDs :many
SELECT
    id,
    name,
    is_active
FROM stores
WHERE id = ANY($1::bigint[]);

-- name: GetStoreDetailByID :one
SELECT
    id,
    name,
    address,
    phone,
    is_active,
    created_at,
    updated_at
FROM stores
WHERE id = $1;

-- name: CheckStoresExistAndActive :one
SELECT
    COUNT(*) as total_count,
    COUNT(CASE WHEN is_active = true THEN 1 END) as active_count
FROM stores
WHERE id = ANY($1::bigint[]);

-- name: CheckStoreNameExists :one
SELECT EXISTS(
    SELECT 1 FROM stores WHERE name = $1
);

-- name: CheckStoreNameExistsExcluding :one
SELECT EXISTS(
    SELECT 1 FROM stores WHERE name = $1 AND id != $2
);

-- name: CheckStoreExistByID :one
SELECT EXISTS(
    SELECT 1 FROM stores WHERE id = $1
);