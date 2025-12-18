package get_orders_by_user

import (
	"context"

	"github.com/4udiwe/big-bob-pizza/order-service/internal/entity"
	"github.com/google/uuid"
)

type OrderService interface {
	GetOrdersByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.Order, int, error)
}
