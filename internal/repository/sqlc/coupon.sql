-- name: CreateCoupon :exec
INSERT INTO coupons (
  id,
  name,
  display_name,
  code,
  discount_rate,
  discount_amount,
  is_active,
  note
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
);

-- name: GetCouponByIDs :many
SELECT id, display_name, code, discount_rate, discount_amount, is_active
FROM coupons
WHERE id = ANY($1::bigint[]);


-- name: CheckCouponExists :one
SELECT EXISTS(
  SELECT 1 FROM coupons
  WHERE id = $1
);

-- name: CheckCouponNameExists :one
SELECT EXISTS(
  SELECT 1 FROM coupons
  WHERE name = $1
);

-- name: CheckCouponNameExistsExcluding :one
SELECT EXISTS(
  SELECT 1 FROM coupons
  WHERE name = $1 AND id != $2
);

-- name: CheckCouponCodeExists :one
SELECT EXISTS(
  SELECT 1 FROM coupons
  WHERE code = $1
);