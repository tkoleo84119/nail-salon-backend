CREATE TABLE IF NOT EXISTS customers (
  id               BIGINT       PRIMARY KEY,
  name             VARCHAR(100) NOT NULL,
  phone            VARCHAR(20)  NOT NULL,
  birthday         DATE         NOT NULL,
  city             VARCHAR(100),
  favorite_shapes  TEXT[],
  favorite_colors  TEXT[],
  favorite_styles  TEXT[],
  is_introvert     BOOLEAN      DEFAULT FALSE,
  referral_source  TEXT[],
  referrer         VARCHAR(100),
  customer_note    TEXT,
  store_note       TEXT,
  level            VARCHAR(20),
  is_blacklisted   BOOLEAN      DEFAULT FALSE,
  created_at       TIMESTAMPTZ  DEFAULT NOW(),
  updated_at       TIMESTAMPTZ  DEFAULT NOW()
);

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

CREATE TABLE IF NOT EXISTS customer_tokens (
  id            BIGINT       PRIMARY KEY,
  customer_id   BIGINT       NOT NULL,
  refresh_token VARCHAR(255) NOT NULL UNIQUE,
  user_agent    TEXT,
  ip_address    INET,
  expired_at    TIMESTAMPTZ  NOT NULL,
  is_revoked    BOOLEAN      DEFAULT FALSE,
  created_at    TIMESTAMPTZ  DEFAULT NOW(),
  updated_at    TIMESTAMPTZ  DEFAULT NOW(),
  FOREIGN KEY (customer_id)  REFERENCES customers(id) ON DELETE CASCADE
);