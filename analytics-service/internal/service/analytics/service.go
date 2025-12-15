package analytics

import (
	"context"
	"time"

	"github.com/4udiwe/big-bob-pizza/analytics-service/internal/entity"
	order_event_repo "github.com/4udiwe/big-bob-pizza/analytics-service/internal/repository/order_event"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// SaveOrderEvent сохраняет событие заказа и обновляет метрики
func (s *Service) SaveOrderEvent(ctx context.Context, event entity.OrderEvent) error {
	logrus.Infof("AnalyticsService.SaveOrderEvent: eventType=%s orderID=%s", event.EventType, event.OrderID)

	if err := s.OrderEventRepo.Save(ctx, event); err != nil {
		logrus.Errorf("AnalyticsService.SaveOrderEvent: failed to save event: %v", err)
		return err
	}

	// Обновляем Prometheus метрики
	s.Metrics.RecordOrderEvent(event.EventType)

	if event.Amount != nil {
		s.Metrics.RecordOrderAmount(*event.Amount, event.EventType)
	}

	logrus.Infof("AnalyticsService.SaveOrderEvent: event saved and metrics updated")
	return nil
}

// GetOrderEvents возвращает все события для заказа
func (s *Service) GetOrderEvents(ctx context.Context, orderID uuid.UUID) ([]entity.OrderEvent, error) {
	logrus.Infof("AnalyticsService.GetOrderEvents: orderID=%s", orderID)
	events, err := s.OrderEventRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		logrus.Errorf("AnalyticsService.GetOrderEvents: error: %v", err)
		return nil, err
	}
	return events, nil
}

// GetStats возвращает статистику за период
func (s *Service) GetStats(ctx context.Context, startDate, endDate time.Time) ([]order_event_repo.OrderStats, error) {
	logrus.Infof("AnalyticsService.GetStats: startDate=%v endDate=%v", startDate, endDate)
	stats, err := s.OrderEventRepo.GetStatsByDateRange(ctx, startDate, endDate)
	if err != nil {
		logrus.Errorf("AnalyticsService.GetStats: error: %v", err)
		return nil, err
	}
	return stats, nil
}

// GetRevenue возвращает выручку за период
func (s *Service) GetRevenue(ctx context.Context, startDate, endDate time.Time) (float64, error) {
	logrus.Infof("AnalyticsService.GetRevenue: startDate=%v endDate=%v", startDate, endDate)
	revenue, err := s.OrderEventRepo.GetTotalRevenue(ctx, startDate, endDate)
	if err != nil {
		logrus.Errorf("AnalyticsService.GetRevenue: error: %v", err)
		return 0, err
	}
	return revenue, nil
}

