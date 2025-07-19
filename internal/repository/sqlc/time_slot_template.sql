-- name: CreateTimeSlotTemplate :one
INSERT INTO time_slot_templates (
    id,
    name,
    note,
    updater,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, NOW(), NOW()
) RETURNING
    id,
    name,
    note,
    updater,
    created_at,
    updated_at;

-- name: BatchCreateTimeSlotTemplateItems :copyfrom
INSERT INTO time_slot_template_items (
    id,
    template_id,
    start_time,
    end_time,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: GetTimeSlotTemplateByID :one
SELECT
    id,
    name,
    note,
    updater,
    created_at,
    updated_at
FROM time_slot_templates
WHERE id = $1;

-- name: GetTimeSlotTemplateItemsByTemplateID :many
SELECT
    id,
    template_id,
    start_time,
    end_time,
    created_at,
    updated_at
FROM time_slot_template_items
WHERE template_id = $1
ORDER BY start_time;

-- name: DeleteTimeSlotTemplate :exec
DELETE FROM time_slot_templates
WHERE id = $1;