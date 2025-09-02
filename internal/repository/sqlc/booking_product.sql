-- name: BulkCreateBookingProducts :copyfrom
INSERT INTO booking_products (
    booking_id,
    product_id,
    created_at
) VALUES (
    $1, $2, $3
);

-- name: GetAllBookingProductIdsByBookingID :many
SELECT product_id FROM booking_products WHERE booking_id = $1;

-- name: CountProductsByIDs :one
SELECT COUNT(*) FROM products
WHERE id = ANY($1::bigint[]) AND store_id = $2;
