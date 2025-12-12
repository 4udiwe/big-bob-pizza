package get_active_orders_by_user

import (
	"context"

	"github.com/4udiwe/big-bob-pizza/order-service/internal/entity"
	"github.com/google/uuid"
)

type OrderService interface {
	GetActiveOrdersByUser(ctx context.Context, userID uuid.UUID) ([]entity.Order, error)
}
