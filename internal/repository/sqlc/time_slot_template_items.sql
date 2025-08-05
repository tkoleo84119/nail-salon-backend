
-- name: CreateTimeSlotTemplateItem :one
INSERT INTO time_slot_template_items (
    id,
    template_id,
    start_time,
    end_time,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, NOW(), NOW()
) RETURNING
    id,
    template_id,
    start_time,
    end_time,
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

-- name: GetTimeSlotTemplateItemsByTemplateIDExcluding :many
SELECT id, template_id, start_time, end_time, created_at, updated_at
FROM time_slot_template_items
WHERE template_id = $1 AND id != $2;

-- name: GetTimeSlotTemplateItemByID :one
SELECT id, template_id, start_time, end_time, created_at, updated_at
FROM time_slot_template_items
WHERE id = $1;

-- name: UpdateTimeSlotTemplateItem :one
UPDATE time_slot_template_items
SET
    start_time = $2,
    end_time = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING id, template_id, start_time, end_time, created_at, updated_at;

-- name: DeleteTimeSlotTemplateItem :exec
DELETE FROM time_slot_template_items
WHERE id = $1;