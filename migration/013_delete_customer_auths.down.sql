CREATE TABLE IF NOT EXISTS customer_auths (
  id           BIGINT       PRIMARY KEY,
  customer_id  BIGINT       NOT NULL,
  provider     VARCHAR(50)  NOT NULL,
  provider_uid VARCHAR(250) NOT NULL,
  other_info   JSONB,
  created_at   TIMESTAMPTZ  DEFAULT NOW(),
  updated_at   TIMESTAMPTZ  DEFAULT NOW(),
  FOREIGN KEY (customer_id)       REFERENCES customers(id) ON DELETE CASCADE,
  CONSTRAINT uq_customer_provider UNIQUE (customer_id, provider),
  CONSTRAINT uq_provider_uid      UNIQUE (provider, provider_uid)
);