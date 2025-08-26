-- name: CreateCustomerCoupon :exec
INSERT INTO customer_coupons (
  id,
  customer_id,
  coupon_id,
  valid_from,
  valid_to
) VALUES (
  $1, $2, $3, $4, $5
);

-- name: GetCustomerCouponPriceInfoByID :one
SELECT
  cc.coupon_id,
  cc.customer_id,
  c.discount_rate,
  c.discount_amount,
  c.is_active,
  cc.valid_to,
  cc.is_used
FROM customer_coupons cc
JOIN coupons c ON cc.coupon_id = c.id
WHERE cc.id = $1;