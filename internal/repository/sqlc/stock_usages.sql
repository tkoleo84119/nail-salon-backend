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
