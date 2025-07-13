-- name: GetAllActiveStores :many
SELECT
    id,
    name
FROM stores
WHERE is_active = true
ORDER BY name;