-- +goose Up
-- +goose StatementBegin
-- ================================
--  Table: dishes
-- ================================
CREATE TABLE IF NOT EXISTS dishes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description VARCHAR(1000),
    price NUMERIC(10,2) NOT NULL,
    category VARCHAR(100) NOT NULL,
    is_available BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_dishes_category ON dishes (category);
CREATE INDEX IF NOT EXISTS idx_dishes_is_available ON dishes (is_available);
CREATE INDEX IF NOT EXISTS idx_dishes_created_at ON dishes (created_at);

-- ================================
--  Table: promotions
-- ================================
CREATE TABLE IF NOT EXISTS promotions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description VARCHAR(1000),
    discount_percent NUMERIC(5,2) NOT NULL,
    dish_ids JSONB NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_promotions_is_active ON promotions (is_active);
CREATE INDEX IF NOT EXISTS idx_promotions_starts_at ON promotions (starts_at);

-- +goose StatementEnd

