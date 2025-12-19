-- +goose Up
-- +goose StatementBegin

-- ================================
--  Table: order_events
--  Хранит события заказов для аналитики
-- ================================
CREATE TABLE order_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL, -- ID события из Kafka (для дедупликации)
    event_type VARCHAR(50) NOT NULL,
    order_id UUID NOT NULL,
    user_id UUID NULL,
    amount NUMERIC(10,2) NULL,
    payment_id UUID NULL,
    reason TEXT NULL, -- причина отмены
    occurred_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_order_events_order_id ON order_events (order_id);
CREATE INDEX idx_order_events_user_id ON order_events (user_id);
CREATE INDEX idx_order_events_event_type ON order_events (event_type);
CREATE INDEX idx_order_events_occurred_at ON order_events (occurred_at);
CREATE UNIQUE INDEX idx_order_events_event_id ON order_events (event_id); -- дедупликация

-- ================================
--  View: order_statistics
--  Статистика по заказам
-- ================================
CREATE VIEW order_statistics AS
SELECT 
    DATE_TRUNC('day', occurred_at) AS date,
    event_type,
    COUNT(*) AS count,
    COUNT(DISTINCT order_id) AS unique_orders,
    COUNT(DISTINCT user_id) AS unique_users,
    SUM(amount) AS total_amount
FROM order_events
WHERE amount IS NOT NULL
GROUP BY DATE_TRUNC('day', occurred_at), event_type;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS order_statistics;
DROP TABLE IF EXISTS order_events;
-- +goose StatementEnd

