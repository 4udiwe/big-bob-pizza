package entity

import (
	"time"

	"github.com/google/uuid"
)

// OrderEvent представляет событие заказа для аналитики
type OrderEvent struct {
	ID         uuid.UUID
	EventID    uuid.UUID // ID события из Kafka
	EventType  string    // order.created, order.paid, order.cancelled, order.completed
	OrderID    uuid.UUID
	UserID     *uuid.UUID // может быть nil для некоторых событий
	Amount     *float64   // может быть nil для некоторых событий
	PaymentID  *uuid.UUID // может быть nil
	Reason     *string    // для отмены
	OccurredAt time.Time
	CreatedAt  time.Time
}

