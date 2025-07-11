CREATE TABLE IF NOT EXISTS suppliers (
  id         BIGINT       PRIMARY KEY,
  name       VARCHAR(100) NOT NULL UNIQUE,
  is_active  BOOLEAN      DEFAULT TRUE,
  created_at TIMESTAMPTZ  DEFAULT NOW(),
  updated_at TIMESTAMPTZ  DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS expenses (
  id            BIGINT        PRIMARY KEY,
  store_id      BIGINT        NOT NULL,
  category      VARCHAR(100),
  supplier_id   BIGINT        NOT NULL,
  amount        NUMERIC(10,2) NOT NULL,
  expense_date  DATE          NOT NULL,
  note          TEXT,
  payer_id      BIGINT,
  is_reimbursed BOOLEAN,
  reimbursed_at TIMESTAMPTZ,
  created_at    TIMESTAMPTZ   DEFAULT NOW(),
  updated_at    TIMESTAMPTZ   DEFAULT NOW(),
  FOREIGN KEY (store_id)    REFERENCES stores(id) ON DELETE CASCADE,
  FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE,
  FOREIGN KEY (payer_id)    REFERENCES staff_users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS expense_items (
  id               BIGINT        PRIMARY KEY,
  expense_id       BIGINT        NOT NULL,
  product_id       BIGINT        NOT NULL,
  quantity         INT           NOT NULL,
  total_price      NUMERIC(10,2) NOT NULL,
  expiration_date  DATE,
  is_arrived       BOOLEAN       DEFAULT FALSE,
  arrival_date     DATE,
  storage_location VARCHAR(100),
  note             TEXT,
  created_at       TIMESTAMPTZ   DEFAULT NOW(),
  updated_at       TIMESTAMPTZ   DEFAULT NOW(),
  FOREIGN KEY (expense_id) REFERENCES expenses(id) ON DELETE CASCADE,
  FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);