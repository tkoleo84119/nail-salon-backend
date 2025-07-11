CREATE TABLE IF NOT EXISTS stores (
  id           BIGINT       PRIMARY KEY,
  name         VARCHAR(100) NOT NULL UNIQUE,
  address      TEXT,
  phone        VARCHAR(20),
  is_active    BOOLEAN      DEFAULT TRUE,
  created_at   TIMESTAMPTZ  DEFAULT NOW(),
  updated_at   TIMESTAMPTZ  DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS staff_users (
  id            BIGINT       PRIMARY KEY,
  username      VARCHAR(50)  NOT NULL UNIQUE,
  email         TEXT         NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  role          VARCHAR(50)  NOT NULL,
  is_active     BOOLEAN      DEFAULT TRUE,
  created_at    TIMESTAMPTZ  DEFAULT NOW(),
  updated_at    TIMESTAMPTZ  DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS staff_user_store_access (
  store_id      BIGINT      NOT NULL,
  staff_user_id BIGINT      NOT NULL,
  created_at    TIMESTAMPTZ DEFAULT NOW(),
  updated_at    TIMESTAMPTZ DEFAULT NOW(),
  PRIMARY KEY (store_id, staff_user_id),
  FOREIGN KEY (store_id)      REFERENCES stores(id) ON DELETE CASCADE,
  FOREIGN KEY (staff_user_id) REFERENCES staff_users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS staff_user_tokens (
  id            BIGINT        PRIMARY KEY,
  staff_user_id BIGINT        NOT NULL,
  refresh_token VARCHAR(255)  NOT NULL UNIQUE,
  user_agent    TEXT,
  ip_address    INET,
  expired_at    TIMESTAMPTZ   NOT NULL,
  is_revoked    BOOLEAN       DEFAULT FALSE,
  created_at    TIMESTAMPTZ   DEFAULT NOW(),
  updated_at    TIMESTAMPTZ   DEFAULT NOW(),
  FOREIGN KEY (staff_user_id) REFERENCES staff_users(id) ON DELETE CASCADE
);
