-- name: CreateCustomer :one
INSERT INTO customers (id, name, phone, birthday, city, favorite_shapes, favorite_colors,
       favorite_styles, is_introvert, referral_source, referrer, customer_note)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING id, name, phone, birthday, city, favorite_shapes, favorite_colors,
       favorite_styles, is_introvert, referral_source, referrer, customer_note,
       store_note, level, is_blacklisted, created_at, updated_at;

-- name: GetCustomerByID :one
SELECT id, name, phone, birthday, city, favorite_shapes, favorite_colors,
       favorite_styles, is_introvert, referral_source, referrer, customer_note,
       store_note, level, is_blacklisted, created_at, updated_at
FROM customers
WHERE id = $1;

-- name: ExistsCustomerByID :one
SELECT EXISTS (SELECT 1 FROM customers WHERE id = $1);