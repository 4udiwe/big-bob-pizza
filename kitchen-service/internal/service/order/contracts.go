package order

import (
	"context"

	"github.com/4udiwe/big-bob-pizza/kitchen-service/internal/entity"
)

type OrderRepository interface {
	StoreOrderStatus(ctx context.Context, orderID string, status string) error
	MarkOrderHanded(ctx context.Context, orderID string) error
	GetActiveOrders(ctx context.Context) ([]entity.KitchenOrder, error)
}
