package post_order

import (
	"context"

	"github.com/4udiwe/big-bob-pizza/order-service/internal/entity"
)

type OrderService interface {
	CreateOrder(ctx context.Context, ord entity.Order) (entity.Order, error)
}
