-- name: BatchCreateExpenseItems :copyfrom
INSERT INTO expense_items (
    id,
    expense_id,
    product_id,
    quantity,
    price,
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

-- name: CreateStoreExpenseItem :exec
INSERT INTO expense_items (
    id,
    expense_id,
    product_id,
    quantity,
    price,
    expiration_date,
    is_arrived,
    arrival_date,
    storage_location,
    note
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
);

-- name: GetStoreExpenseItemsByExpenseID :many
SELECT
    ei.id,
    ei.product_id,
    p.name AS product_name,
    ei.quantity,
    ei.price,
    ei.expiration_date,
    ei.is_arrived,
    ei.arrival_date,
    ei.storage_location,
    ei.note
FROM expense_items ei
LEFT JOIN products p ON ei.product_id = p.id
WHERE ei.expense_id = $1;

-- name: GetStoreExpenseItemByID :one
SELECT
    id,
    expense_id,
    product_id,
    quantity,
    price,
    expiration_date,
    is_arrived,
    arrival_date,
    storage_location,
    note
FROM expense_items
WHERE id = $1
AND expense_id = $2;

-- name: CheckExpenseItemsExistsByExpenseID :one
SELECT EXISTS(SELECT 1 FROM expense_items WHERE expense_id = $1);

-- name: CheckAllExpenseItemsAreArrived :one
SELECT NOT EXISTS(
    SELECT 1 FROM expense_items WHERE expense_id = $1 AND is_arrived = false
) AS all_arrived;

-- name: DeleteStoreExpenseItem :exec
DELETE FROM expense_items
WHERE id = $1
AND expense_id = $2;