package analytics

import (
	"context"
	"time"

	"github.com/4udiwe/big-bob-pizza/analytics-service/internal/entity"
	"github.com/google/uuid"
)

import (
	order_event_repo "github.com/4udiwe/big-bob-pizza/analytics-service/internal/repository/order_event"
)

type OrderEventRepo interface {
	Save(ctx context.Context, event entity.OrderEvent) error
	GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]entity.OrderEvent, error)
	GetStatsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]order_event_repo.OrderStats, error)
	GetTotalRevenue(ctx context.Context, startDate, endDate time.Time) (float64, error)
}

type Service struct {
	OrderEventRepo OrderEventRepo
	Metrics        *Metrics
}

func NewService(orderEventRepo OrderEventRepo, metrics *Metrics) *Service {
	return &Service{
		OrderEventRepo: orderEventRepo,
		Metrics:        metrics,
	}
}

