-- name: CreateProduct :exec
INSERT INTO products (
    id,
    store_id,
    name,
    brand_id,
    category_id,
    current_stock,
    safety_stock,
    unit,
    storage_location,
    note
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
);

-- name: GetProductByID :one
SELECT
    id,
    store_id,
    name,
    brand_id,
    category_id,
    current_stock,
    safety_stock,
    unit,
    storage_location,
    note,
    created_at,
    updated_at
FROM products
WHERE id = $1;

-- name: CheckProductNameBrandExistsInStore :one
SELECT EXISTS(
    SELECT 1 FROM products
    WHERE store_id = $1 AND name = $2 AND brand_id = $3
);

-- name: CheckProductNameBrandExistsInStoreExcluding :one
SELECT EXISTS(
    SELECT 1 FROM products
    WHERE store_id = $1 AND name = $2 AND brand_id = $3 AND id != $4
);