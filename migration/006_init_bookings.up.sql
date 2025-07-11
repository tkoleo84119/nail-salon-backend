CREATE TABLE IF NOT EXISTS services (
    id               BIGINT        PRIMARY KEY,
    name             VARCHAR(150)  NOT NULL UNIQUE,
    price            NUMERIC(10,2) NOT NULL,
    duration_minutes INT           NOT NULL,
    is_addon         BOOLEAN       DEFAULT FALSE,
    is_visible       BOOLEAN       DEFAULT TRUE,
    is_active        BOOLEAN       DEFAULT TRUE,
    note             TEXT,
    created_at       TIMESTAMPTZ   DEFAULT NOW(),
    updated_at       TIMESTAMPTZ   DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS coupons (
    id               BIGINT       PRIMARY KEY,
    name             VARCHAR(100) NOT NULL,
    code             VARCHAR(50)  NOT NULL UNIQUE,
    discount_rate    NUMERIC(3,2),
    discount_amount  NUMERIC(10,2),
    is_active        BOOLEAN      DEFAULT TRUE,
    created_at       TIMESTAMPTZ  DEFAULT NOW(),
    updated_at       TIMESTAMPTZ  DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS bookings (
    id              BIGINT       PRIMARY KEY,
    store_id        BIGINT       NOT NULL,
    customer_id     BIGINT       NOT NULL,
    stylist_id      BIGINT       NOT NULL,
    time_slot_id    BIGINT       NOT NULL,
    is_chat_enabled BOOLEAN      DEFAULT TRUE,
    no_show         BOOLEAN      DEFAULT FALSE,
    actual_duration INT,
    note            TEXT,
    used_products   TEXT[],
    status          VARCHAR(30)  NOT NULL,
    created_at      TIMESTAMPTZ  DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  DEFAULT NOW(),
    FOREIGN KEY (store_id)     REFERENCES stores(id) ON DELETE CASCADE,
    FOREIGN KEY (customer_id)  REFERENCES customers(id) ON DELETE CASCADE,
    FOREIGN KEY (stylist_id)   REFERENCES stylists(id) ON DELETE CASCADE,
    FOREIGN KEY (time_slot_id) REFERENCES time_slots(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS checkouts (
    id              BIGINT        PRIMARY KEY,
    booking_id      BIGINT        NOT NULL,
    total_amount    NUMERIC(12,2) NOT NULL,
    final_amount    NUMERIC(12,2) NOT NULL,
    paid_amount     NUMERIC(12,2) NOT NULL,
    payment_method  VARCHAR(50)   NOT NULL,
    coupon_id       BIGINT,
    checkout_user   BIGINT,
    created_at      TIMESTAMPTZ   DEFAULT NOW(),
    updated_at      TIMESTAMPTZ   DEFAULT NOW(),
    FOREIGN KEY (booking_id)    REFERENCES bookings(id) ON DELETE CASCADE,
    FOREIGN KEY (coupon_id)     REFERENCES coupons(id) ON DELETE CASCADE,
    FOREIGN KEY (checkout_user) REFERENCES staff_users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS booking_details (
    id               BIGINT      PRIMARY KEY,
    booking_id       BIGINT      NOT NULL,
    service_id       BIGINT      NOT NULL,
    price            NUMERIC(10,2),
    discount_rate    NUMERIC(3,2),
    discount_amount  NUMERIC(10,2),
    created_at       TIMESTAMPTZ DEFAULT NOW(),
    updated_at       TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (booking_id) REFERENCES bookings(id) ON DELETE CASCADE,
    FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS customer_coupons (
    id          BIGINT      PRIMARY KEY,
    customer_id BIGINT      NOT NULL,
    coupon_id   BIGINT      NOT NULL,
    valid_from  TIMESTAMPTZ NOT NULL,
    valid_to    TIMESTAMPTZ,
    is_used     BOOLEAN     DEFAULT FALSE,
    used_at     TIMESTAMPTZ,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (customer_id)     REFERENCES customers(id) ON DELETE CASCADE,
    FOREIGN KEY (coupon_id)       REFERENCES coupons(id) ON DELETE CASCADE,
    CONSTRAINT uq_customer_coupon UNIQUE (customer_id, coupon_id)
);