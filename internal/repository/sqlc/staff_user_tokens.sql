-- name: CreateStaffUserToken :one
INSERT INTO staff_user_tokens (
    id,
    staff_user_id,
    refresh_token,
    user_agent,
    ip_address,
    expired_at
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING id, created_at;

-- name: GetValidStaffUserToken :one
SELECT id, staff_user_id, refresh_token, user_agent, ip_address, expired_at,
       is_revoked, created_at, updated_at
FROM staff_user_tokens
WHERE refresh_token = $1 AND expired_at > NOW() AND is_revoked = false;