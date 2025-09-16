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
SELECT id, staff_user_id
FROM staff_user_tokens
WHERE refresh_token = $1 AND expired_at > NOW() AND is_revoked = false;

-- name: RevokeStaffUserToken :exec
UPDATE staff_user_tokens
SET is_revoked = true
WHERE refresh_token = $1;

-- name: CountExpiredOrRevokedStaffUserTokens :one
SELECT COUNT(*) FROM staff_user_tokens
WHERE is_revoked = true OR expired_at < NOW();

-- name: DeleteStaffUserTokensBatch :exec
DELETE FROM staff_user_tokens
WHERE id IN (
    SELECT id
    FROM staff_user_tokens
    WHERE is_revoked = true OR expired_at < NOW()
    LIMIT $1
);