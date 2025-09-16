-- name: CreateCustomerToken :one
INSERT INTO customer_tokens (id, customer_id, refresh_token, user_agent, ip_address, expired_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, customer_id, refresh_token, user_agent, ip_address, expired_at, is_revoked, created_at, updated_at;

-- name: GetValidCustomerToken :one
SELECT id, customer_id
FROM customer_tokens
WHERE refresh_token = $1 AND expired_at > NOW() AND is_revoked = false;

-- name: RevokeCustomerToken :exec
UPDATE customer_tokens
SET is_revoked = true
WHERE refresh_token = $1;

-- name: CountExpiredOrRevokedCustomerTokens :one
SELECT COUNT(*) FROM customer_tokens
WHERE is_revoked = true OR expired_at < NOW();

-- name: DeleteCustomerTokensBatch :exec
DELETE FROM customer_tokens
WHERE id IN (
  SELECT id
  FROM customer_tokens
  WHERE is_revoked = true OR expired_at < NOW()
  LIMIT $1
);