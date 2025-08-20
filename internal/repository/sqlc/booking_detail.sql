
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
    price,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: GetBookingDetailsByBookingID :many
SELECT
    bd.id,
    bd.booking_id,
    bd.service_id,
    srv.name as service_name,
    bd.price,
    bd.created_at,
    srv.is_addon
FROM booking_details bd
JOIN services srv ON bd.service_id = srv.id
WHERE bd.booking_id = $1
ORDER BY srv.is_addon ASC, srv.name ASC;

-- name: GetBookingDetailsByBookingIDs :many
SELECT
    bd.booking_id,
    bd.service_id,
    srv.name AS service_name,
    srv.is_addon
FROM booking_details bd
JOIN services srv ON bd.service_id = srv.id
WHERE bd.booking_id = ANY($1::bigint[])
ORDER BY bd.booking_id ASC, srv.is_addon ASC, srv.name;

-- name: DeleteBookingDetailsByBookingID :exec
DELETE FROM booking_details
WHERE booking_id = $1;