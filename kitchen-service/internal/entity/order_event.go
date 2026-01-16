package entity

import (
	"time"

	"github.com/google/uuid"
)

// OrderEvent представляет событие заказа из order-service
type OrderEvent struct {
	ID         uuid.UUID
	EventID    uuid.UUID
	EventType  string
	OrderID    uuid.UUID
	UserID     *uuid.UUID
	Amount     *float64
	PaymentID  *uuid.UUID
	Reason     *string
	OccurredAt time.Time
	CreatedAt  time.Time
}
