-- name: CreateSupplier :one
INSERT INTO suppliers (id, name)
VALUES ($1, $2)
RETURNING id;

-- name: CheckSupplierNameExists :one
SELECT EXISTS(SELECT 1 FROM suppliers WHERE name = $1);

-- name: CheckSupplierExistsByID :one
SELECT EXISTS(SELECT 1 FROM suppliers WHERE id = $1);

-- name: CheckSupplierNameExistsExcluding :one
SELECT EXISTS(SELECT 1 FROM suppliers WHERE name = $1 AND id != $2);