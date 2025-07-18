
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
