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