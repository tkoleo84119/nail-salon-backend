-- name: CreateAccountTransaction :one
INSERT INTO account_transactions (
    id,
    account_id,
    transaction_date,
    type,
    amount,
    balance,
    note
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING id;

-- name: GetAccountTransactionByID :one
SELECT id, account_id, transaction_date, type, amount, balance, note
FROM account_transactions
WHERE id = $1;

-- name: GetAccountTransactionCurrentBalance :one
SELECT COALESCE(
(SELECT balance
    FROM account_transactions
    WHERE account_id = $1
    ORDER BY transaction_date DESC, created_at DESC
    LIMIT 1),
    0
)::int as balance;
