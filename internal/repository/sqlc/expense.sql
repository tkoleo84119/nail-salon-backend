-- name: CheckSupplierExists :one
SELECT EXISTS(SELECT 1 FROM suppliers WHERE id = $1 AND is_active = true);

-- name: CreateExpense :one
INSERT INTO expenses (
    id,
    store_id,
    category,
    supplier_id,
    amount,
    expense_date,
    note,
    payer_id,
    is_reimbursed
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING id;

-- name: GetStoreExpenseByID :one
SELECT
    e.id,
    e.supplier_id,
    COALESCE(s.name, '') AS supplier_name,
    e.payer_id,
    COALESCE(su.username, '') AS payer_name,
    e.category,
    e.amount,
    e.expense_date,
    e.note,
    e.is_reimbursed,
    e.reimbursed_at,
    e.created_at,
    e.updated_at
FROM expenses e
LEFT JOIN suppliers s ON e.supplier_id = s.id
LEFT JOIN staff_users su ON e.payer_id = su.id
WHERE e.id = $1 AND e.store_id = $2;

-- name: GetStoreExpenseItemsByExpenseID :many
SELECT
    ei.id,
    ei.product_id,
    p.name AS product_name,
    ei.quantity,
    ei.total_price,
    ei.expiration_date,
    ei.is_arrived,
    ei.arrival_date,
    ei.storage_location,
    ei.note
FROM expense_items ei
LEFT JOIN products p ON ei.product_id = p.id
WHERE ei.expense_id = $1;