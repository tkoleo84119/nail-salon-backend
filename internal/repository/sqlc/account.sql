-- name: CreateAccount :exec
INSERT INTO accounts (
    id,
    store_id,
    name,
    note
) VALUES (
    $1, $2, $3, $4
);
