CREATE TABLE IF NOT EXISTS accounts (
  id         BIGINT       PRIMARY KEY,
  store_id   BIGINT       NOT NULL,
  name       VARCHAR(100) NOT NULL,
  note       TEXT,
  is_active  BOOLEAN      DEFAULT TRUE,
  created_at TIMESTAMPTZ  DEFAULT NOW(),
  updated_at TIMESTAMPTZ  DEFAULT NOW(),
  FOREIGN KEY (store_id)           REFERENCES stores(id) ON DELETE CASCADE,
  CONSTRAINT uq_account_store_name UNIQUE (store_id, name)
);

CREATE TABLE IF NOT EXISTS account_transactions (
  id               BIGINT        PRIMARY KEY,
  account_id       BIGINT        NOT NULL,
  transaction_date TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
  type             VARCHAR(50)   NOT NULL,
  amount           NUMERIC(12,2) NOT NULL,
  balance          NUMERIC(12,2) NOT NULL,
  note             TEXT,
  created_at       TIMESTAMPTZ   DEFAULT NOW(),
  updated_at       TIMESTAMPTZ   DEFAULT NOW(),
  FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);