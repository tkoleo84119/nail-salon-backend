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