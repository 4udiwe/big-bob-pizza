package payment

import (
	"context"

	"github.com/4udiwe/big-bob-pizza/order-service/pkg/transactor"
	"github.com/4udiwe/big-bob-pizza/payment-service/internal/entity"
	"github.com/google/uuid"
)

type PaymentRepo interface {
	Create(ctx context.Context, payment entity.Payment) (entity.Payment, error)
	GetByID(ctx context.Context, paymentID uuid.UUID) (entity.Payment, error)
	GetByOrderID(ctx context.Context, orderID uuid.UUID) (entity.Payment, error)
	UpdateStatus(ctx context.Context, paymentID uuid.UUID, status entity.PaymentStatus, failureReason *string) error
	GetAllPayments(ctx context.Context, limit, offset int, status *entity.PaymentStatusName, userID *uuid.UUID) ([]entity.PaymentWithUser, int, error)
	GetPaymentsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.PaymentWithUser, int, error)
}


type OrderCacheRepo interface {
	GetByOrderID(ctx context.Context, orderID uuid.UUID) (entity.OrderInfo, error)
	Delete(ctx context.Context, orderID uuid.UUID) error
}

type OutboxRepo interface {
	Create(ctx context.Context, ev entity.OutboxEvent) error
}

type Service struct {
	PaymentRepo    PaymentRepo
	OrderCacheRepo OrderCacheRepo
	OutboxRepo     OutboxRepo
	TxManager      transactor.Transactor
}

func NewService(
	paymentRepo PaymentRepo,
	orderCacheRepo OrderCacheRepo,
	outboxRepo OutboxRepo,
	txManager transactor.Transactor,
) *Service {
	return &Service{
		PaymentRepo:    paymentRepo,
		OrderCacheRepo: orderCacheRepo,
		OutboxRepo:     outboxRepo,
		TxManager:      txManager,
	}
}
