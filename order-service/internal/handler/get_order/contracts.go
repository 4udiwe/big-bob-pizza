package get_order

import (
	"context"

	"github.com/4udiwe/big-bob-pizza/order-service/internal/entity"
	"github.com/google/uuid"
)

type OrderService interface {
	GetOrderByID(ctx context.Context, orderID uuid.UUID) (entity.Order, error)
}
