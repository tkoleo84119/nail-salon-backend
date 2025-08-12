
-- name: CreateBookingDetail :one
INSERT INTO booking_details (
    id,
    booking_id,
    service_id,
    price
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: CreateBookingDetails :copyfrom
INSERT INTO booking_details (
    id,
    booking_id,
    service_id,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5
);