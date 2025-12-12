-- +goose Up
-- +goose StatementBegin
-- ================================
--  Lookup table: payment_status
-- ================================
CREATE TABLE payment_status (
    id SMALLSERIAL PRIMARY KEY,
    name VARCHAR(32) NOT NULL UNIQUE
);

INSERT INTO payment_status (name) VALUES
('pending'), ('completed'), ('failed');

-- ================================
--  Lookup table: outbox_status
-- ================================
CREATE TABLE outbox_status (
    id SMALLSERIAL PRIMARY KEY,
    name VARCHAR(32) NOT NULL UNIQUE
);

INSERT INTO outbox_status (name) VALUES
('pending'), ('processed'), ('failed');

-- ================================
--  Table: payments
-- ================================
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL,
    amount NUMERIC(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'RUB',
    status_id SMALLINT NOT NULL REFERENCES payment_status(id) DEFAULT 1,
    failure_reason TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payment_order_id ON payments (order_id);
CREATE INDEX idx_payment_status_id ON payments (status_id);
CREATE INDEX idx_payment_created_at ON payments (created_at);

-- ================================
--  Table: outbox
-- ================================
CREATE TABLE outbox (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aggregate_type VARCHAR(64) NOT NULL,
    aggregate_id UUID NOT NULL,
    event_type VARCHAR(128) NOT NULL,
    payload JSONB NOT NULL,
    status_id SMALLINT NOT NULL REFERENCES outbox_status(id) DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_outbox_status_id ON outbox (status_id);
CREATE INDEX idx_outbox_created_at ON outbox (created_at);

-- ================================
--  Table: order_cache (для хранения информации о заказах, доступных для оплаты)
-- ================================
CREATE TABLE order_cache (
    order_id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    total_price NUMERIC(10,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_order_cache_expires_at ON order_cache (expires_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS order_cache;
DROP TABLE IF EXISTS outbox;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS outbox_status;
DROP TABLE IF EXISTS payment_status;

-- +goose StatementEnd

