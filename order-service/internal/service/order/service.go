package order

import (
	"context"

	"github.com/4udiwe/avito-pvz/pkg/transactor"
	"github.com/4udiwe/big-bob-pizza/order-service/internal/entity"
	"github.com/google/uuid"
)

type Service struct {
	orderRepo OrderRepository
	txManager transactor.Transactor
}

func New(orderRepo OrderRepository, txManager transactor.Transactor) *Service {
	return &Service{
		orderRepo: orderRepo,
		txManager: txManager,
	}
}

func (s *Service) CreateOrder(ctx context.Context, customerID uuid.UUID, totalAmount float64, currency string) (entity.Product, error)

func (s *Service) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status entity.OrderStatus) error

func (s *Service) GetOrderByID(ctx context.Context, orderID uuid.UUID) (entity.Product, error)

func (s *Service) GetAllOrders(ctx context.Context) ([]entity.Product, error)
