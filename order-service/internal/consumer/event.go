package consumer

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

// Все типы событий, которые могут потребляться сервисом
const (
	PaymentSuccess EventType = "payment.success"
	PaymentFailed  EventType = "payment.failed"

	KitchenAccepted        EventType = "kitchen.accepted"
	KitchenReady           EventType = "kitchen.ready"
	KitchenHandedToCourier EventType = "kitchen.handedToCourier"

	DeliveryCompleted EventType = "delivery.completed"
)

// Тип для обработки входящего события
type IncomingEvent struct {
	Type       EventType
	OccurredAt time.Time
	Payload    Payload
}

// Тип данных входящего события. Перечисленны все поля, которые могут быть в событиию. 
// (omitempty опускает поле, если его нет)
type Payload struct {
	PaymentID   uuid.UUID `json:"paymentId"`
	OrderID     uuid.UUID `json:"orderId,omitempty"`
	Amount      float64   `json:"amount,omitempty"`
	Reason      string    `json:"reason,omitempty"`
	DeliveryID  uuid.UUID `json:"deliveryId,omitempty"`
	DeliveredAt time.Time `json:"deliveredAt,omitempty"`
}
