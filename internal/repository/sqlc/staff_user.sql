-- name: GetActiveStaffUserByUsername :one
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

-- name: CheckStaffUserExistsByUsername :one
SELECT EXISTS(
    SELECT 1 FROM staff_users
    WHERE username = $1
) as exists;

-- name: CreateStaffUser :one
INSERT INTO staff_users (
    id,
    username,
    email,
    password_hash,
    role,
    is_active,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, true, NOW(), NOW()
) RETURNING
    id,
    username,
    email,
    role,
    is_active,
    created_at,
    updated_at;

-- name: UpdateStaffUserPassword :one
UPDATE staff_users
SET password_hash = $2, updated_at = NOW()
WHERE id = $1
RETURNING id;
