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