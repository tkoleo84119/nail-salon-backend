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

-- name: CheckProductNameBrandExistsInStore :one
SELECT EXISTS(
    SELECT 1 FROM products
    WHERE store_id = $1 AND name = $2 AND brand_id = $3
);
