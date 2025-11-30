package events

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	PaymentSuccess EventType = "payment.success"
	PaymentFailed  EventType = "payment.failed"

	KitchenAccepted        EventType = "kitchen.accepted"
	KitchenReady           EventType = "kitchen.ready"
	KitchenHandedToCourier EventType = "kitchen.handedToCourier"

	DeliveryCompleted EventType = "delivery.completed"
)

type IncomingEvent struct {
	Type       EventType
	OccurredAt time.Time
	Payload    Payload
}

type Payload struct {
	PaymentID   uuid.UUID `json:"paymentId"`
	OrderID     uuid.UUID `json:"orderId,omitempty"`
	Amount      float64   `json:"amount,omitempty"`
	Reason      string    `json:"reason,omitempty"`
	DeliveryID  uuid.UUID `json:"deliveryId,omitempty"`
	DeliveredAt time.Time `json:"deliveredAt,omitempty"`
}
