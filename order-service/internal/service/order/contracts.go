package order

import (
	"context"

	"github.com/4udiwe/big-bob-pizza/order-service/internal/entity"
	"github.com/google/uuid"
)

//go:generate go tool mockgen -source=contracts.go -destination=mocks/repo_mock.go -package=mocks

type OrderRepository interface {
	CreateOrder(ctx context.Context, customerID uuid.UUID, totalAmount float64, currency string) (entity.Product, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status entity.OrderStatus) error
	GetOrderByID(ctx context.Context, orderID uuid.UUID) (entity.Product, error)
	GetAllOrders(ctx context.Context) ([]entity.Product, error)
}
