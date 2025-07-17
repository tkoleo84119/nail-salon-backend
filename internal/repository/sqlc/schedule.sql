-- name: CreateSchedule :one
INSERT INTO schedules (
    id,
    store_id,
    stylist_id,
    work_date,
    note,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, NOW(), NOW()
) RETURNING
    id,
    store_id,
    stylist_id,
    work_date,
    note,
    created_at,
    updated_at;

-- name: BatchCreateSchedules :copyfrom
INSERT INTO schedules (
    id,
    store_id,
    stylist_id,
    work_date,
    note,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- name: CheckScheduleExists :one
SELECT EXISTS(
    SELECT 1 FROM schedules
    WHERE store_id = $1 AND stylist_id = $2 AND work_date = $3
) as exists;

-- name: GetSchedulesByStoreAndStylist :many
SELECT
    id,
    store_id,
    stylist_id,
    work_date,
    note,
    created_at,
    updated_at
FROM schedules
WHERE store_id = $1 AND stylist_id = $2
ORDER BY work_date;

-- name: GetScheduleByID :one
SELECT
    id,
    store_id,
    stylist_id,
    work_date,
    note,
    created_at,
    updated_at
FROM schedules
WHERE id = $1;

-- name: GetSchedulesWithTimeSlotsByIDs :many
SELECT
    s.id,
    s.store_id,
    s.stylist_id,
    s.work_date,
    s.note,
    s.created_at,
    s.updated_at,
    t.id as time_slot_id,
    t.start_time,
    t.end_time,
    t.is_available
FROM schedules s
LEFT JOIN time_slots t ON s.id = t.schedule_id
WHERE s.id = ANY($1::bigint[])
ORDER BY s.work_date;

-- name: DeleteSchedulesByIDs :exec
DELETE FROM schedules
WHERE id = ANY($1::bigint[]);
