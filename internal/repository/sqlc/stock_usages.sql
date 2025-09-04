-- name: CreateStockUsage :exec
INSERT INTO stock_usages (
    id,
    product_id,
    quantity,
    expiration,
    usage_started
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: GetStockUsageByID :one
SELECT
    id,
    product_id,
    quantity,
    is_in_use,
    expiration,
    usage_started,
    usage_ended_at,
    created_at,
    updated_at
FROM stock_usages
WHERE id = $1;

-- name: UpdateStockUsageFinish :exec
UPDATE stock_usages
SET
    is_in_use = false,
    usage_ended_at = $2,
    updated_at = NOW()
WHERE id = $1;
