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

-- name: GetScheduleByID :one
SELECT
    id,
    store_id,
    stylist_id,
    work_date
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

-- name: GetScheduleWithTimeSlotsByID :many
SELECT
    s.id,
    s.work_date,
    s.note,
    t.id as time_slot_id,
    t.start_time,
    t.end_time,
    t.is_available
FROM schedules s
LEFT JOIN time_slots t ON s.id = t.schedule_id
WHERE s.id = $1
ORDER BY t.start_time;

-- name: GetAvailableSchedules :many
SELECT s.id, s.work_date
FROM schedules s
JOIN time_slots ts ON s.id = ts.schedule_id
WHERE s.store_id = $1
  AND s.stylist_id = $2
  AND s.work_date BETWEEN $3 AND $4
  AND ts.is_available = true
GROUP BY s.id, s.work_date
ORDER BY s.work_date ASC;

-- name: DeleteSchedulesByIDs :exec
DELETE FROM schedules
WHERE id = ANY($1::bigint[]);

-- name: CheckScheduleDateExists :one
SELECT EXISTS(
    SELECT 1 FROM schedules
    WHERE store_id = $1 AND stylist_id = $2 AND work_date = $3
) as exists;

-- name: CheckScheduleExistsByID :one
SELECT EXISTS(
    SELECT 1 FROM schedules
    WHERE id = $1
) as exists;

-- name: CheckScheduleCanUpdateDate :one
SELECT NOT EXISTS(
    SELECT 1
    FROM time_slots ts
    WHERE ts.schedule_id = $1
    AND (
        ts.is_available = false
        OR (
            ts.is_available = true
            AND EXISTS (
                SELECT 1
                FROM bookings b
                WHERE b.time_slot_id = ts.id
                AND b.status IN ('CANCELLED', 'NO_SHOW')
            )
        )
    )
) as can_update;