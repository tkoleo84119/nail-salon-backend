CREATE TABLE IF NOT EXISTS customer_terms_acceptance (
    id BIGINT PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    terms_version VARCHAR(50) NOT NULL,
    accepted_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
);