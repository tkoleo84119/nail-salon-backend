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
JOIN stylists st ON b.stylist_id = st.id
JOIN time_slots ts ON b.time_slot_id = ts.id
JOIN schedules sch ON ts.schedule_id = sch.id
WHERE b.id = $1;

-- name: CancelBooking :one
UPDATE bookings
SET status = $2, cancel_reason = $3, updated_at = NOW()
WHERE id = $1 AND customer_id = $4
RETURNING id, status, cancel_reason, updated_at;

-- name: UpdateBookingByStaff :one
UPDATE bookings
SET
    time_slot_id = COALESCE($2, time_slot_id),
    is_chat_enabled = COALESCE($3, is_chat_enabled),
    note = COALESCE($4, note),
    updated_at = NOW()
WHERE id = $1
RETURNING id;
