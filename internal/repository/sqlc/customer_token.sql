-- name: CreateCustomerToken :one
INSERT INTO customer_tokens (id, customer_id, refresh_token, user_agent, ip_address, expired_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, customer_id, refresh_token, user_agent, ip_address, expired_at, is_revoked, created_at, updated_at;

-- name: GetValidCustomerToken :one
SELECT id, customer_id
FROM customer_tokens
WHERE refresh_token = $1 AND expired_at > NOW() AND is_revoked = false;