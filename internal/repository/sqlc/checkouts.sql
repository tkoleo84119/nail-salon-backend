-- name: CreateCheckout :exec
INSERT INTO checkouts (
  id,
  booking_id,
  total_amount,
  final_amount,
  paid_amount,
  payment_method,
  coupon_id,
  checkout_user
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
);

-- name: GetCheckoutByBookingID :one
SELECT
  ck.id,
  ck.total_amount,
  ck.final_amount,
  ck.paid_amount,
  ck.payment_method,
  ck.coupon_id,
  c.name as coupon_name,
  c.code as coupon_code,
  su.username as checkout_user
FROM checkouts ck
LEFT JOIN coupons c ON c.id = ck.coupon_id
LEFT JOIN staff_users su ON su.id = ck.checkout_user
WHERE ck.booking_id = $1;