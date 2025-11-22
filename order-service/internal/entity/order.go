package entity

import "github.com/google/uuid"

type OrderStatus string

const (
	OrderStatusCreated    OrderStatus = "created"
	OrderStatusPaid       OrderStatus = "paid"
	OrderStatusPreparing  OrderStatus = "preparing"
	OrderStatusDelivering OrderStatus = "delivering"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

type Product struct {
	ID          uuid.UUID   `db:"id"`
	CustomerID  uuid.UUID   `db:"customer_id"`
	Status      OrderStatus `db:"status"`
	TotalAmount float64     `db:"total_amount"`
	Currency    string      `db:"currency"`
	PaymentID   *uuid.UUID  `db:"payment_id"`
	DeliveryID  *uuid.UUID  `db:"delivery_id"`
	CreatedAt   string      `db:"created_at"`
	UpdatedAt   string      `db:"updated_at"`
}
