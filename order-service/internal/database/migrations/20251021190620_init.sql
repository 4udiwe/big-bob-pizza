-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ================================
--  Lookup table: order_status
-- ================================
CREATE TABLE order_status (
    id SMALLSERIAL PRIMARY KEY,
    name VARCHAR(32) NOT NULL UNIQUE
);

INSERT INTO order_status (name) VALUES
('created'), ('paid'), ('preparing'), ('delivering'), ('completed'), ('cancelled');

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
--  Table: order
-- ================================
CREATE TABLE order (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_id UUID NOT NULL,
    status_id SMALLINT NOT NULL REFERENCES order_status(id),
    total_amount NUMERIC(10,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'RUB',
    payment_id UUID NULL,
    delivery_id UUID NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_order_customer_id ON order (customer_id);
CREATE INDEX idx_order_status_id ON order (status_id);
CREATE INDEX idx_order_updated_at ON order (updated_at);

-- ================================
--  Table: order_item
-- ================================
CREATE TABLE order_item (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES order(id) ON DELETE CASCADE,
    product_id UUID NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    product_price NUMERIC(10,2) NOT NULL,
    amount INT NOT NULL CHECK (amount > 0),
    total_price NUMERIC(10,2) GENERATED ALWAYS AS (product_price * amount) STORED,
    notes TEXT NULL
);

CREATE INDEX idx_order_items_order_id ON order_item (order_id);

-- ================================
--  Table: outbox
-- ================================
CREATE TABLE outbox (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    aggregate_type VARCHAR(64) NOT NULL,
    aggregate_id UUID NOT NULL,
    event_type VARCHAR(128) NOT NULL,
    payload JSONB NOT NULL,
    status_id SMALLINT NOT NULL REFERENCES outbox_statuses(id) DEFAULT 1, -- pending
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_outbox_status_id ON outbox (status_id);
CREATE INDEX idx_outbox_created_at ON outbox (created_at);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS outbox;
DROP TABLE IF EXISTS order_item;
DROP TABLE IF EXISTS order;
DROP TABLE IF EXISTS outbox_status;
DROP TABLE IF EXISTS order_status;

-- +goose StatementEnd
