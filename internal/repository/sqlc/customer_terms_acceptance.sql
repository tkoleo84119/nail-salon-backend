-- name: CreateCustomerTermsAcceptance :exec
INSERT INTO customer_terms_acceptance (id, customer_id, terms_version, accepted_at)
VALUES ($1, $2, $3, $4);

-- name: CheckCustomerTermsExistsByCustomerIDAndVersion :one
SELECT EXISTS(
  SELECT 1 FROM customer_terms_acceptance
  WHERE customer_id = $1 AND terms_version = $2
) as exists;

-- name: GetCustomerTermsAcceptanceByCustomerIDAndVersion :one
SELECT id, customer_id, terms_version, accepted_at
FROM customer_terms_acceptance
WHERE customer_id = $1 AND terms_version = $2;