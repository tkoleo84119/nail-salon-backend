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
    c.line_uid as customer_line_uid,
    c.name as customer_name,
    c.phone as customer_phone,
    b.stylist_id,
    st.name as stylist_name,
    b.time_slot_id,
    ts.start_time,
    ts.end_time,
    sch.work_date,
    b.is_chat_enabled,
    b.note,
    b.actual_duration,
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

-- name: GetBookingInfoWithDateByID :one
SELECT
    b.status,
    b.customer_id,
    b.store_id,
    sch.work_date,
    ts.start_time
FROM bookings b
JOIN time_slots ts ON b.time_slot_id = ts.id
JOIN schedules sch ON ts.schedule_id = sch.id
WHERE b.id = $1;

-- name: UpdateBookingsStatus :exec
UPDATE bookings
SET status = $2, updated_at = NOW()
WHERE id = ANY($1::bigint[]);

-- name: UpdateBookingActualDuration :exec
UPDATE bookings
SET actual_duration = $2, updated_at = NOW()
WHERE id = $1;

-- name: CancelBooking :one
UPDATE bookings
SET status = $2, cancel_reason = $3, updated_at = NOW()
WHERE id = $1
RETURNING id;

-- name: GetStylistPerformanceGroupByStore :many
SELECT
    b.store_id,
    s.name as store_name,
    COUNT(*) as total_bookings,
    SUM(CASE WHEN b.status = 'COMPLETED' THEN 1 ELSE 0 END) as completed_bookings,
    SUM(CASE WHEN b.status = 'CANCELLED' THEN 1 ELSE 0 END) as cancelled_bookings,
    SUM(CASE WHEN b.status = 'NO_SHOW' THEN 1 ELSE 0 END) as no_show_bookings,
    COALESCE(SUM(CASE WHEN c.payment_method = 'LINE_PAY' AND b.status = 'COMPLETED' THEN COALESCE(c.final_amount, 0) ELSE 0 END), 0)::numeric(12,2) as line_pay_revenue,
    COALESCE(SUM(CASE WHEN c.payment_method = 'CASH' AND b.status = 'COMPLETED' THEN COALESCE(c.final_amount, 0) ELSE 0 END), 0)::numeric(12,2) as cash_revenue,
    COALESCE(SUM(CASE WHEN b.status = 'COMPLETED' THEN COALESCE(c.paid_amount, 0) ELSE 0 END), 0)::numeric(12,2) as total_paid_amount,
    SUM(COALESCE(b.actual_duration, 0)) as total_service_time
FROM bookings b
INNER JOIN stores s ON b.store_id = s.id
INNER JOIN time_slots ts ON b.time_slot_id = ts.id
INNER JOIN schedules sch ON ts.schedule_id = sch.id
LEFT JOIN checkouts c ON b.id = c.booking_id
WHERE b.stylist_id = $1
    AND b.status != 'SCHEDULE'
    AND sch.work_date BETWEEN $2 AND $3
GROUP BY b.store_id, s.name
ORDER BY b.store_id;

-- name: GetStorePerformanceGroupByStylist :many
SELECT
    b.stylist_id,
    st.name as stylist_name,
    COUNT(*) as total_bookings,
    SUM(CASE WHEN b.status = 'COMPLETED' THEN 1 ELSE 0 END) as completed_bookings,
    SUM(CASE WHEN b.status = 'CANCELLED' THEN 1 ELSE 0 END) as cancelled_bookings,
    SUM(CASE WHEN b.status = 'NO_SHOW' THEN 1 ELSE 0 END) as no_show_bookings,
    COALESCE(SUM(CASE WHEN c.payment_method = 'LINE_PAY' AND b.status = 'COMPLETED' THEN COALESCE(c.final_amount, 0) ELSE 0 END), 0)::numeric(12,2) as line_pay_revenue,
    COALESCE(SUM(CASE WHEN c.payment_method = 'CASH' AND b.status = 'COMPLETED' THEN COALESCE(c.final_amount, 0) ELSE 0 END), 0)::numeric(12,2) as cash_revenue,
    COALESCE(SUM(CASE WHEN b.status = 'COMPLETED' THEN COALESCE(c.paid_amount, 0) ELSE 0 END), 0)::numeric(12,2) as total_paid_amount,
    SUM(COALESCE(b.actual_duration, 0)) as total_service_time
FROM bookings b
INNER JOIN stores s ON b.store_id = s.id
INNER JOIN stylists st ON b.stylist_id = st.id
INNER JOIN time_slots ts ON b.time_slot_id = ts.id
INNER JOIN schedules sch ON ts.schedule_id = sch.id
LEFT JOIN checkouts c ON b.id = c.booking_id
WHERE b.store_id = $1
    AND b.status != 'SCHEDULE'
    AND sch.work_date BETWEEN $2 AND $3
GROUP BY b.stylist_id, st.name
ORDER BY b.stylist_id;

-- name: GetTomorrowBookingsForReminder :many
SELECT
    b.id,
    b.store_id,
    s.name as store_name,
    s.address as store_address,
    b.customer_id,
    c.line_uid as customer_line_uid,
    c.name as customer_name,
    b.time_slot_id,
    ts.start_time,
    ts.end_time,
    sch.work_date,
    b.status
FROM bookings b
JOIN stores s ON b.store_id = s.id
JOIN customers c ON b.customer_id = c.id
JOIN time_slots ts ON b.time_slot_id = ts.id
JOIN schedules sch ON ts.schedule_id = sch.id
WHERE sch.work_date = $1
    AND b.status = 'SCHEDULED'
ORDER BY s.id, ts.start_time;