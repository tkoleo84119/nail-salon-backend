-- name: CreateCustomer :exec
INSERT INTO customers (id, line_uid, line_name, email, name, phone, birthday, city, favorite_shapes, favorite_colors,
      favorite_styles, is_introvert, referral_source, referrer, customer_note, level)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16);

-- name: GetCustomerByID :one
SELECT id, name, line_uid, line_name, phone, birthday, email, city, favorite_shapes, favorite_colors,
      favorite_styles, is_introvert, referral_source, referrer, customer_note,
      store_note, level, is_blacklisted, last_visit_at, created_at, updated_at
FROM customers
WHERE id = $1;

-- name: GetCustomerByIDs :many
SELECT id, name, line_name, phone
FROM customers
WHERE id = ANY($1::bigint[]);

-- name: GetCustomerByLineUid :one
SELECT id, line_name, name
FROM customers
WHERE line_uid = $1;

-- name: UpdateCustomerLineName :exec
UPDATE customers
SET line_name = $2
WHERE id = $1;

-- name: UpdateCustomerLastVisitAt :exec
UPDATE customers
SET last_visit_at = NOW()
WHERE id = $1;

-- name: CheckCustomerExistsByID :one
SELECT EXISTS (SELECT 1 FROM customers WHERE id = $1);

-- name: CheckCustomerExistsByLineUid :one
SELECT EXISTS (SELECT 1 FROM customers WHERE line_uid = $1);