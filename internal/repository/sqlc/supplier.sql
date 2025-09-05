-- name: CreateSupplier :one
INSERT INTO suppliers (id, name)
VALUES ($1, $2)
RETURNING id;

-- name: CheckSupplierNameExists :one
SELECT EXISTS(SELECT 1 FROM suppliers WHERE name = $1);
