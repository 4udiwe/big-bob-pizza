package entity

import "time"

type KitchenOrderStatus string

const (
	KitchenOrderStatusPending   KitchenOrderStatus = "PENDING"   // заказ получен кухней
	KitchenOrderStatusCooking   KitchenOrderStatus = "COOKING"   // готовится
	KitchenOrderStatusReady     KitchenOrderStatus = "READY"     // готов
	KitchenOrderStatusCancelled KitchenOrderStatus = "CANCELLED" // отменён
)

type KitchenOrder struct {
	ID        string
	OrderID   string // ID заказа из order-service
	Status    KitchenOrderStatus
	Items     []KitchenItem
	CreatedAt time.Time
	UpdatedAt time.Time
}

type KitchenItem struct {
	DishID   string
	Name     string
	Quantity int
}
