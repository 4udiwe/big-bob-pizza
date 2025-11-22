-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ================================
--  Table: order
-- ================================
CREATE TYPE order_status AS ENUM (
    'created',
    'paid',
    'preparing',
    'delivering',
    'completed',
    'cancelled'
);

CREATE TABLE order (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_id UUID NOT NULL,
    status order_status NOT NULL DEFAULT 'created',
    total_amount NUMERIC(10,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'RUB',
    payment_id UUID NULL,
    delivery_id UUID NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_order_customer_id ON order (customer_id);
CREATE INDEX idx_order_status ON order (status);

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

CREATE INDEX idx_order_item_order_id ON order_item (order_id);

-- ================================
--  Table: outbox
-- ================================
CREATE TYPE outbox_status AS ENUM (
    'pending',
    'processed',
    'failed'
);

CREATE TABLE outbox (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    aggregate_type VARCHAR(64) NOT NULL,
    aggregate_id UUID NOT NULL,
    event_type VARCHAR(128) NOT NULL,
    payload JSONB NOT NULL,
    status outbox_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_outbox_status ON outbox (status);
CREATE INDEX idx_outbox_created_at ON outbox (created_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS outbox;
DROP TABLE IF EXISTS order_item;
DROP TABLE IF EXISTS order;

DROP TYPE IF EXISTS outbox_status;
DROP TYPE IF EXISTS order_status;
-- +goose StatementEnd
