CREATE TABLE IF NOT EXISTS brands (
  id         BIGINT       PRIMARY KEY,
  name       VARCHAR(100) NOT NULL UNIQUE,
  is_active  BOOLEAN      DEFAULT TRUE,
  created_at TIMESTAMPTZ  DEFAULT NOW(),
  updated_at TIMESTAMPTZ  DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS product_categories (
  id         BIGINT       PRIMARY KEY,
  name       VARCHAR(100) NOT NULL UNIQUE,
  is_active  BOOLEAN      DEFAULT TRUE,
  created_at TIMESTAMPTZ  DEFAULT NOW(),
  updated_at TIMESTAMPTZ  DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS products (
  id               BIGINT       PRIMARY KEY,
  store_id         BIGINT       NOT NULL,
  name             VARCHAR(200) NOT NULL,
  brand_id         BIGINT       NOT NULL,
  category_id      BIGINT       NOT NULL,
  current_stock    INT          NOT NULL DEFAULT 0,
  safety_stock     INT          DEFAULT -1,
  unit             VARCHAR(50),
  storage_location VARCHAR(100),
  note             TEXT,
  created_at       TIMESTAMPTZ  DEFAULT NOW(),
  updated_at       TIMESTAMPTZ  DEFAULT NOW(),
  FOREIGN KEY (store_id)                 REFERENCES stores(id) ON DELETE CASCADE,
  FOREIGN KEY (brand_id)                 REFERENCES brands(id) ON DELETE CASCADE,
  FOREIGN KEY (category_id)              REFERENCES product_categories(id) ON DELETE CASCADE,
  CONSTRAINT uq_product_store_brand_name UNIQUE (store_id, brand_id, name)
);

CREATE TABLE IF NOT EXISTS stock_usages (
  id             BIGINT      PRIMARY KEY,
  product_id     BIGINT      NOT NULL,
  quantity       INT         NOT NULL,
  is_in_use      BOOLEAN     DEFAULT TRUE,
  expiration     DATE        NOT NULL,
  usage_started  DATE        NOT NULL,
  usage_ended_at DATE,
  created_at     TIMESTAMPTZ DEFAULT NOW(),
  updated_at     TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);
