CREATE TABLE IF NOT EXISTS stylists (
  id             BIGINT       PRIMARY KEY,
  staff_user_id  BIGINT       UNIQUE,
  name           VARCHAR(100),
  good_at_shapes TEXT[],
  good_at_colors TEXT[],
  good_at_styles TEXT[],
  is_introvert   BOOLEAN      DEFAULT FALSE,
  created_at     TIMESTAMPTZ  DEFAULT NOW(),
  updated_at     TIMESTAMPTZ  DEFAULT NOW(),
  FOREIGN KEY (staff_user_id) REFERENCES staff_users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS schedules (
  id         BIGINT      PRIMARY KEY,
  store_id   BIGINT      NOT NULL,
  stylist_id BIGINT      NOT NULL,
  work_date  DATE        NOT NULL,
  note       TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (store_id)   REFERENCES stores(id) ON DELETE CASCADE,
  FOREIGN KEY (stylist_id) REFERENCES stylists(id) ON DELETE CASCADE,
  CONSTRAINT uq_schedule   UNIQUE (store_id, stylist_id, work_date)
);

CREATE TABLE IF NOT EXISTS time_slots (
  id           BIGINT      PRIMARY KEY,
  schedule_id  BIGINT      NOT NULL,
  start_time   TIME        NOT NULL,
  end_time     TIME        NOT NULL,
  is_available BOOLEAN     DEFAULT TRUE,
  created_at   TIMESTAMPTZ DEFAULT NOW(),
  updated_at   TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (schedule_id) REFERENCES schedules(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS time_slot_templates (
  id          BIGINT       PRIMARY KEY,
  name        VARCHAR(100) NOT NULL,
  note        TEXT,
  updater     BIGINT,
  created_at  TIMESTAMPTZ  DEFAULT NOW(),
  updated_at  TIMESTAMPTZ  DEFAULT NOW(),
  FOREIGN KEY (updater) REFERENCES staff_users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS time_slot_template_items (
  id          BIGINT      PRIMARY KEY,
  template_id BIGINT      NOT NULL,
  start_time  TIME        NOT NULL,
  end_time    TIME        NOT NULL,
  created_at  TIMESTAMPTZ DEFAULT NOW(),
  updated_at  TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (template_id) REFERENCES time_slot_templates(id) ON DELETE CASCADE
);
