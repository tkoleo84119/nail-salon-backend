CREATE INDEX IF NOT EXISTS idx_staff_user_tokens_on_revoked_expired ON staff_user_tokens (is_revoked, expired_at);
CREATE INDEX IF NOT EXISTS idx_customer_tokens_on_revoked_expired ON customer_tokens (is_revoked, expired_at);

CREATE INDEX IF NOT EXISTS idx_staff_user_tokens_on_staff_user_id ON staff_user_tokens (staff_user_id);
CREATE INDEX IF NOT EXISTS idx_customer_tokens_on_customer_id ON customer_tokens (customer_id);