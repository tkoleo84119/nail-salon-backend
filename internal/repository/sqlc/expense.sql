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
