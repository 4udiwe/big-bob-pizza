package entity

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID          uuid.UUID
	CustomerID  uuid.UUID
	Status      OrderStatus
	TotalAmount float64
	Currency    string
	PaymentID   *uuid.UUID
	DeliveryID  *uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Items       []OrderItem
}
