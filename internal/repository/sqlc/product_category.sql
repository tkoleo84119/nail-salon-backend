-- name: CreateProductCategory :one
INSERT INTO product_categories (id, name)
VALUES ($1, $2)
RETURNING id;

-- name: CheckProductCategoryNameExists :one
SELECT EXISTS(SELECT 1 FROM product_categories WHERE name = $1);

-- name: CheckProductCategoryNameExistsExcludeSelf :one
SELECT EXISTS(SELECT 1 FROM product_categories WHERE name = $1 AND id != $2);

-- name: CheckProductCategoryExistByID :one
SELECT EXISTS(SELECT 1 FROM product_categories WHERE id = $1);
