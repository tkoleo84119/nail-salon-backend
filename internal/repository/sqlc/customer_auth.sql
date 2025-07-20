-- name: CreateCustomerAuth :one
INSERT INTO customer_auths (id, customer_id, provider, provider_uid, other_info)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, customer_id, provider, provider_uid, other_info, created_at, updated_at;

-- name: GetCustomerAuthByProviderUid :one
SELECT ca.id, ca.customer_id, ca.provider, ca.provider_uid, ca.other_info,
      ca.created_at, ca.updated_at,
      c.id as customer_id, c.name as customer_name, c.phone as customer_phone
FROM customer_auths ca
JOIN customers c ON ca.customer_id = c.id
WHERE ca.provider = $1 AND ca.provider_uid = $2;