package order

import (
	"context"
	"time"

	"github.com/4udiwe/big-bob-pizza/order-service/internal/entity"
	"github.com/google/uuid"
)

type OrderRepo interface {
	// Inserts only order data (without order items).
	// Receives entity with order data.
	Create(ctx context.Context, order entity.Order) (entity.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status entity.StatusName, time time.Time) error
	UpdateOrderPayment(ctx context.Context, orderID, paymentID uuid.UUID, time time.Time) error
	UpdateOrderDelivery(ctx context.Context, orderID, deliveryID uuid.UUID, time time.Time) error
	GetOrderByID(ctx context.Context, orderID uuid.UUID) (entity.Order, error)
	// Return page of found orders sorted by creation time, total items, and error.
	GetAllOrders(ctx context.Context, limit, offset int) (orders []entity.Order, total int, err error)
	GetOrdersByUserID(ctx context.Context, userID uuid.UUID) (orders []entity.Order, err error)
}

type ItemsRepo interface {
	// Inserts items for order.
	// Returns filled items (with TotalPrice and ID).
	InsertItems(ctx context.Context, orderID uuid.UUID, items []entity.OrderItem) ([]entity.OrderItem, error)
}

type OutboxRepo interface {
	Create(ctx context.Context, ev entity.OutboxEvent) error
}

type CacheRepo interface {
	Save(ctx context.Context, ord *entity.Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Order, error)
	AddToActive(ctx context.Context, ord *entity.Order) error
	RemoveFromActive(ctx context.Context, id uuid.UUID) error
	AddUserActive(ctx context.Context, userID uuid.UUID, ordID uuid.UUID) error
	RemoveUserActive(ctx context.Context, userID uuid.UUID, ordID uuid.UUID) error
	AddToStatus(ctx context.Context, status string, ordID uuid.UUID) error
	RemoveFromStatus(ctx context.Context, status string, ordID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetActiveOrders(ctx context.Context) ([]string, error)
	GetUserActiveOrders(ctx context.Context, userID uuid.UUID) ([]string, error)
	GetByStatus(ctx context.Context, status string) ([]string, error)
}
