-- name: BatchCreateExpenseItems :copyfrom
INSERT INTO expense_items (
    id,
    expense_id,
    product_id,
    quantity,
    total_price,
    expiration_date,
    is_arrived,
    arrival_date,
    storage_location,
    note,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
);