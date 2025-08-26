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