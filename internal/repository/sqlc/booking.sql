-- name: CreateBooking :one
INSERT INTO bookings (
    id,
    store_id,
    customer_id,
    stylist_id,
    time_slot_id,
    is_chat_enabled,
    note,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;
