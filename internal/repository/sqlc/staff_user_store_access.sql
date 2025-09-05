-- name: CreateStaffUserStoreAccess :exec
INSERT INTO staff_user_store_access (
    store_id,
    staff_user_id,
    created_at,
    updated_at
) VALUES (
    $1, $2, NOW(), NOW()
);

-- name: GetAllActiveStoreAccessByStaffId :many
SELECT
    sa.store_id,
    s.name as store_name
FROM staff_user_store_access sa
JOIN stores s ON sa.store_id = s.id
WHERE sa.staff_user_id = $1 AND s.is_active = true;

-- name: BatchCreateStaffUserStoreAccess :copyfrom
INSERT INTO staff_user_store_access (
    store_id,
    staff_user_id,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4
);

-- name: CheckStoreAccessExists :one
SELECT EXISTS(
    SELECT 1 FROM staff_user_store_access
    WHERE staff_user_id = $1 AND store_id = $2
) as exists;

-- name: DeleteStaffUserStoreAccess :exec
DELETE FROM staff_user_store_access
WHERE staff_user_id = $1 AND store_id = ANY($2::bigint[]);

-- name: CheckStaffHasStoreAccess :one
SELECT EXISTS(
    SELECT 1
    FROM staff_user_store_access susa
    JOIN staff_users su ON susa.staff_user_id = su.id
    WHERE susa.staff_user_id = $1
    AND susa.store_id = $2
    AND su.is_active = true
);