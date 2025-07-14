-- name: GetStaffUserByUsername :one
SELECT
    id,
    username,
    email,
    password_hash,
    role,
    is_active,
    created_at,
    updated_at
FROM staff_users
WHERE username = $1 AND is_active = true;

-- name: GetStaffUserByID :one
SELECT
    id,
    username,
    email,
    password_hash,
    role,
    is_active,
    created_at,
    updated_at
FROM staff_users
WHERE id = $1;





