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

-- name: CheckCouponNameExists :one
SELECT EXISTS(
  SELECT 1 FROM coupons
  WHERE name = $1
);

-- name: CheckCouponCodeExists :one
SELECT EXISTS(
  SELECT 1 FROM coupons
  WHERE code = $1
);