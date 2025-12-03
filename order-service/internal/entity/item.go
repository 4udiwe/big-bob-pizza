package entity

import "github.com/google/uuid"

type OrderItem struct {
	ID           uuid.UUID
	ProductID    uuid.UUID
	ProductName  string
	ProductPrice float64
	Amount       int
	TotalPrice   float64
	Notes        string
}
