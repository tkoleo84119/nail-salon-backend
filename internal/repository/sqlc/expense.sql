-- name: CheckSupplierExists :one
SELECT EXISTS(SELECT 1 FROM suppliers WHERE id = $1 AND is_active = true);

-- name: CreateExpense :one
INSERT INTO expenses (
    id,
    store_id,
    category,
    supplier_id,
    amount,
    other_fee,
    expense_date,
    note,
    payer_id,
    is_reimbursed,
    updater
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
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
    e.other_fee,
    e.expense_date,
    e.note,
    e.is_reimbursed,
    e.reimbursed_at,
    COALESCE(su2.username, '') AS updater,
    e.created_at,
    e.updated_at
FROM expenses e
LEFT JOIN suppliers s ON e.supplier_id = s.id
LEFT JOIN staff_users su ON e.payer_id = su.id
LEFT JOIN staff_users su2 ON e.updater = su2.id
WHERE e.id = $1 AND e.store_id = $2;

-- name: UpdateStoreExpenseAmount :exec
UPDATE expenses SET amount = $1, updater = $2, updated_at = NOW() WHERE id = $3;

-- name: GetExpenseReportSummary :one
SELECT
    COUNT(*) as total_count,
    COALESCE(SUM(amount + COALESCE(other_fee, 0)), 0)::numeric(12,2) as total_amount,
    COALESCE(SUM(CASE WHEN payer_id IS NOT NULL THEN amount + COALESCE(other_fee, 0) ELSE 0 END), 0)::numeric(12,2) as advance_amount,
    COALESCE(SUM(CASE WHEN payer_id IS NOT NULL AND is_reimbursed = true THEN amount + COALESCE(other_fee, 0) ELSE 0 END), 0)::numeric(12,2) as reimbursed_amount
FROM expenses
WHERE store_id = $1
    AND expense_date BETWEEN $2 AND $3;

-- name: GetExpenseReportByCategory :many
SELECT
    COALESCE(category, '未分類') as category,
    COUNT(*) as count,
    COALESCE(SUM(amount + COALESCE(other_fee, 0)), 0)::numeric(12,2) as amount
FROM expenses
WHERE store_id = $1
    AND expense_date BETWEEN $2 AND $3
GROUP BY category;

-- name: GetExpenseReportBySupplier :many
SELECT
    e.supplier_id,
    s.name as supplier_name,
    COUNT(*) as count,
    COALESCE(SUM(e.amount + COALESCE(e.other_fee, 0)), 0)::numeric(12,2) as amount
FROM expenses e
LEFT JOIN suppliers s ON e.supplier_id = s.id
WHERE e.store_id = $1
    AND e.expense_date BETWEEN $2 AND $3
    AND e.supplier_id IS NOT NULL
GROUP BY e.supplier_id, s.name;

-- name: GetExpenseReportByPayer :many
SELECT
    e.payer_id,
    su.username as payer_name,
    COUNT(*) as advance_count,
    COALESCE(SUM(e.amount + COALESCE(e.other_fee, 0)), 0)::numeric(12,2) as advance_amount,
    COALESCE(SUM(CASE WHEN e.is_reimbursed = true THEN e.amount + COALESCE(e.other_fee, 0) ELSE 0 END), 0)::numeric(12,2) as reimbursed_amount
FROM expenses e
LEFT JOIN staff_users su ON e.payer_id = su.id
WHERE e.store_id = $1
    AND e.expense_date BETWEEN $2 AND $3
    AND e.payer_id IS NOT NULL
GROUP BY e.payer_id, su.username;