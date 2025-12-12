package get_all_orders

import (
	"context"

	"github.com/4udiwe/big-bob-pizza/order-service/internal/entity"
)

type OrderService interface {
	GetAllOrders(ctx context.Context, limit, offset int) ([]entity.Order, int, error)
}
