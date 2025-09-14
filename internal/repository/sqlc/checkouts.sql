-- name: BulkCreateCheckout :copyfrom
INSERT INTO checkouts (
  id,
  booking_id,
  total_amount,
  final_amount,
  paid_amount,
  payment_method,
  coupon_id,
  checkout_user,
  created_at,
  updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
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
  c.display_name as coupon_display_name,
  c.code as coupon_code,
  su.username as checkout_user,
  ck.created_at
FROM checkouts ck
LEFT JOIN coupons c ON c.id = ck.coupon_id
LEFT JOIN staff_users su ON su.id = ck.checkout_user
WHERE ck.booking_id = $1;