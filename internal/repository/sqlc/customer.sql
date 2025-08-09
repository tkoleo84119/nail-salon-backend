-- name: CreateCustomer :exec
INSERT INTO customers (id, line_uid, line_name, name, phone, birthday, city, favorite_shapes, favorite_colors,
      favorite_styles, is_introvert, referral_source, referrer, customer_note, level)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15);

-- name: GetCustomerByID :one
SELECT id, name, line_name, phone, birthday, email, city, favorite_shapes, favorite_colors,
      favorite_styles, is_introvert, referral_source, referrer, customer_note,
      store_note, level, is_blacklisted, last_visit_at, created_at, updated_at
FROM customers
WHERE id = $1;

-- name: GetCustomerByLineUid :one
SELECT id
FROM customers
WHERE line_uid = $1;

-- name: ExistsCustomerByID :one
SELECT EXISTS (SELECT 1 FROM customers WHERE id = $1);

-- name: CheckCustomerExistsByLineUid :one
SELECT EXISTS (SELECT 1 FROM customers WHERE line_uid = $1);