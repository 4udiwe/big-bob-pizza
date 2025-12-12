package post_payment

import (
	"context"

	"github.com/4udiwe/big-bob-pizza/payment-service/internal/entity"
	"github.com/google/uuid"
)

type PaymentService interface {
	ProcessPayment(ctx context.Context, orderID uuid.UUID, amount float64) (entity.Payment, error)
}
