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

-- name: GetBookingByID :one
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

-- name: GetBookingDetailsByBookingID :many
SELECT
    bd.id,
    bd.booking_id,
    bd.service_id,
    srv.name as service_name,
    bd.price,
    bd.created_at
FROM booking_details bd
JOIN services srv ON bd.service_id = srv.id
WHERE bd.booking_id = $1
ORDER BY srv.is_addon ASC, srv.name ASC;

-- name: DeleteBookingDetailsByBookingID :exec
DELETE FROM booking_details
WHERE booking_id = $1;
