-- name: GetAllActiveStores :many
SELECT
    id,
    name
FROM stores
WHERE is_active = true
ORDER BY name;

-- name: GetStoresByIDs :many
SELECT
    id,
    name,
    is_active
FROM stores
WHERE id = ANY($1::bigint[]);

-- name: CheckStoresExistAndActive :one
SELECT
    COUNT(*) as total_count,
    COUNT(CASE WHEN is_active = true THEN 1 END) as active_count
FROM stores
WHERE id = ANY($1::bigint[]);