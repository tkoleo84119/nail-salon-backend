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

-- name: CheckStylistExistsByStaffUserID :one
SELECT EXISTS(
    SELECT 1 FROM stylists 
    WHERE staff_user_id = $1
) as exists;

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