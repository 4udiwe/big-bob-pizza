package entity

import (
	"time"

	"github.com/google/uuid"
)

// OrderInfo - информация о заказе, полученная из события order.created
type OrderInfo struct {
	OrderID    uuid.UUID
	UserID     uuid.UUID
	TotalPrice float64
	Items      []OrderItem
	CreatedAt  time.Time
}

type OrderItem struct {
	DishID   uuid.UUID
	Quantity int
}

