package item_repository

import (
	"github.com/4udiwe/big-bob-pizza/order-service/internal/entity"
	"github.com/google/uuid"
)

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
