package order_repository

import (
	"time"

	"github.com/4udiwe/big-bob-pizza/order-service/internal/entity"
	"github.com/google/uuid"
)

type RowOrder struct {
	ID          uuid.UUID  `db:"id"`
	CustomerID  uuid.UUID  `db:"customer_id"`
	StatusID    int        `db:"status_id"`
	StatusName  string     `db:"status_name"`
	TotalAmount float64    `db:"total_amount"`
	Currency    string     `db:"currency"`
	PaymentID   *uuid.UUID `db:"payment_id"`
	DeliveryID  *uuid.UUID `db:"delivery_id"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

func (r *RowOrder) ToEntity() entity.Order {
	return entity.Order{
		ID:          r.ID,
		CustomerID:  r.CustomerID,
		Status:      entity.OrderStatus{ID: r.StatusID, Name: r.StatusName},
		TotalAmount: r.TotalAmount,
		Currency:    r.Currency,
		PaymentID:   r.PaymentID,
		DeliveryID:  r.DeliveryID,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

type RowItem struct {
	ID           uuid.UUID `db:"id"`
	OrderID      uuid.UUID `db:"order_id"`
	ProductID    uuid.UUID `db:"product_id"`
	ProductName  string    `db:"product_name"`
	ProductPrice float64   `db:"product_price"`
	Amount       int       `db:"amount"`
	TotalPrice   float64   `db:"total_price"`
	Notes        string    `db:"notes"`
}

func (r *RowItem) ToEntity() entity.OrderItem {
	return entity.OrderItem{
		ID:           r.ID,
		ProductID:    r.ProductID,
		ProductName:  r.ProductName,
		ProductPrice: r.ProductPrice,
		Amount:       r.Amount,
		TotalPrice:   r.TotalPrice,
		Notes:        r.Notes,
	}
}
