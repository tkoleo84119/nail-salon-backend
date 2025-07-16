-- name: GetStaffUserStoreAccess :many
SELECT
    sa.store_id,
    s.name as store_name
FROM staff_user_store_access sa
JOIN stores s ON sa.store_id = s.id
WHERE sa.staff_user_id = $1 AND s.is_active = true;

-- name: CreateStaffUserStoreAccess :exec
INSERT INTO staff_user_store_access (
    store_id,
    staff_user_id,
    created_at,
    updated_at
) VALUES (
    $1, $2, NOW(), NOW()
);

-- name: BatchCreateStaffUserStoreAccess :exec
INSERT INTO staff_user_store_access (
    store_id,
    staff_user_id,
    created_at,
    updated_at
)
SELECT 
    unnest($1::bigint[]) as store_id,
    $2 as staff_user_id,
    NOW() as created_at,
    NOW() as updated_at;

-- name: CheckStoreAccessExists :one
SELECT EXISTS(
    SELECT 1 FROM staff_user_store_access 
    WHERE staff_user_id = $1 AND store_id = $2
) as exists;

-- name: DeleteStaffUserStoreAccess :exec
DELETE FROM staff_user_store_access 
WHERE staff_user_id = $1 AND store_id = ANY($2::bigint[]);