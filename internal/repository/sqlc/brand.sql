-- name: CreateBrand :one
INSERT INTO brands (id, name)
VALUES ($1, $2)
RETURNING id, name, is_active, created_at, updated_at;

-- name: CheckBrandNameExists :one
SELECT EXISTS(SELECT 1 FROM brands WHERE name = $1);

-- name: CheckBrandNameExistsExcludeSelf :one
SELECT EXISTS(SELECT 1 FROM brands WHERE name = $1 AND id != $2);

-- name: CheckBrandExistByID :one
SELECT EXISTS(SELECT 1 FROM brands WHERE id = $1);