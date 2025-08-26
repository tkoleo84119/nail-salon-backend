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

-- name: GetBookingDetailByID :one
SELECT
    b.id,
    b.store_id,
    s.name as store_name,
    b.customer_id,
    c.name as customer_name,
    b.stylist_id,
    st.name as stylist_name,
    b.time_slot_id,
    ts.start_time,
    ts.end_time,
    sch.work_date,
    b.is_chat_enabled,
    b.note,
    b.status,
    b.created_at,
    b.updated_at
FROM bookings b
JOIN stores s ON b.store_id = s.id
JOIN customers c ON b.customer_id = c.id
JOIN stylists st ON b.stylist_id = st.id
JOIN time_slots ts ON b.time_slot_id = ts.id
JOIN schedules sch ON ts.schedule_id = sch.id
WHERE b.id = $1;

-- name: GetBookingInfoByID :one
SELECT status, customer_id, store_id FROM bookings WHERE id = $1;

-- name: UpdateBookingStatus :exec
UPDATE bookings
SET status = $2, updated_at = NOW()
WHERE id = $1;

-- name: CancelBooking :one
UPDATE bookings
SET status = $2, cancel_reason = $3, updated_at = NOW()
WHERE id = $1
RETURNING id;
