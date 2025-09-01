-- name: CreateBrand :one
INSERT INTO brands (id, name)
VALUES ($1, $2)
RETURNING id, name, is_active, created_at, updated_at;

-- name: CheckBrandNameExists :one
SELECT EXISTS(SELECT 1 FROM brands WHERE name = $1);
