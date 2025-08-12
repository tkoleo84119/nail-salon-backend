-- name: CreateStylist :one
INSERT INTO stylists (
    id,
    staff_user_id,
    name,
    good_at_shapes,
    good_at_colors,
    good_at_styles,
    is_introvert,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, NOW(), NOW()
) RETURNING
    id,
    staff_user_id,
    name,
    good_at_shapes,
    good_at_colors,
    good_at_styles,
    is_introvert,
    created_at,
    updated_at;

-- name: GetStylistByStaffUserID :one
SELECT
    id,
    staff_user_id,
    name,
    good_at_shapes,
    good_at_colors,
    good_at_styles,
    is_introvert,
    created_at,
    updated_at
FROM stylists
WHERE staff_user_id = $1;

-- name: GetStylistByID :one
SELECT
    id,
    staff_user_id,
    name,
    good_at_shapes,
    good_at_colors,
    good_at_styles,
    is_introvert,
    created_at,
    updated_at
FROM stylists
WHERE id = $1;

-- name: CheckStylistExistAndActive :one
SELECT EXISTS(
    SELECT 1
    FROM stylists
    JOIN staff_users ON stylists.staff_user_id = staff_users.id
    WHERE stylists.id = $1
    AND staff_users.is_active = true
);

-- name: GetActiveStylistNameByID :one
SELECT
    name
FROM stylists
JOIN staff_users ON stylists.staff_user_id = staff_users.id
WHERE stylists.id = $1
AND staff_users.is_active = true;