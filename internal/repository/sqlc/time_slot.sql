-- name: CreateTimeSlot :one
INSERT INTO time_slots (
    id,
    schedule_id,
    start_time,
    end_time,
    is_available,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, true, NOW(), NOW()
) RETURNING
    id,
    schedule_id,
    start_time,
    end_time,
    is_available,
    created_at,
    updated_at;


-- name: BatchCreateTimeSlots :copyfrom
INSERT INTO time_slots (
    id,
    schedule_id,
    start_time,
    end_time,
    is_available,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- name: GetTimeSlotsByScheduleID :many
SELECT
    id,
    schedule_id,
    start_time,
    end_time,
    is_available,
    created_at,
    updated_at
FROM time_slots
WHERE schedule_id = $1
ORDER BY start_time;

-- name: GetTimeSlotByID :one
SELECT
    id,
    schedule_id,
    start_time,
    end_time,
    is_available,
    created_at,
    updated_at
FROM time_slots
WHERE id = $1;

-- name: DeleteTimeSlotByID :exec
DELETE FROM time_slots
WHERE id = $1;

-- name: DeleteTimeSlotsByScheduleIDs :exec
DELETE FROM time_slots
WHERE schedule_id = ANY($1::bigint[]);

-- name: CheckTimeSlotOverlap :one
SELECT EXISTS(
    SELECT 1 FROM time_slots
    WHERE schedule_id = $1
    AND start_time < $3
    AND end_time > $2
) AS has_overlap;

-- name: CheckTimeSlotOverlapExcluding :one
SELECT EXISTS(
    SELECT 1 FROM time_slots
    WHERE schedule_id = $1
    AND id != $2
    AND start_time < $4
    AND end_time > $3
) AS has_overlap;

-- name: GetAvailableTimeSlotsByScheduleID :many
SELECT
    ts.id,
    ts.schedule_id,
    ts.start_time,
    ts.end_time,
    ts.is_available,
    ts.created_at,
    ts.updated_at
FROM time_slots ts
LEFT JOIN bookings b ON ts.id = b.time_slot_id AND b.status != 'CANCELLED'
WHERE ts.schedule_id = $1
  AND ts.is_available = true
  AND b.id IS NULL
ORDER BY ts.start_time;