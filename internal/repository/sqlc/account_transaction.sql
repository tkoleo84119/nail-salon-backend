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

-- name: GetAccountTransactionCurrentBalance :one
SELECT COALESCE(
(SELECT balance
    FROM account_transactions
    WHERE account_id = $1
    ORDER BY transaction_date DESC, created_at DESC
    LIMIT 1),
    0
)::int as balance;
