-- name: CreateAccount :exec
INSERT INTO accounts (
    id,
    store_id,
    name,
    note
) VALUES (
    $1, $2, $3, $4
);

-- name: GetAccountByID :one
SELECT id, store_id, name, note, is_active FROM accounts WHERE id = $1;
