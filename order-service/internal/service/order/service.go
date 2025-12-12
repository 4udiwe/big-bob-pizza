package order

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"

	"github.com/4udiwe/big-bob-pizza/order-service/internal/entity"
	"github.com/4udiwe/big-bob-pizza/order-service/internal/repository"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/transactor"
)

type Service struct {
	OrderRepo  OrderRepo
	ItemsRepo  ItemsRepo
	OutboxRepo OutboxRepo
	CacheRepo  CacheRepo
	TxManager  transactor.Transactor
}

func NewService(
	orderRepo OrderRepo,
	itemsRepo ItemsRepo,
	outboxRepo OutboxRepo,
	cacheRepo CacheRepo,
	txManager transactor.Transactor,
) *Service {
	return &Service{
		OrderRepo:  orderRepo,
		ItemsRepo:  itemsRepo,
		OutboxRepo: outboxRepo,
		CacheRepo:  cacheRepo,
		TxManager:  txManager,
	}
}

func (s *Service) CreateOrder(ctx context.Context, ord entity.Order) (entity.Order, error) {
	log.Infof("OrderService.CreateOrder: creating order for customer %s", ord.CustomerID)

	var created entity.Order

	err := s.TxManager.WithinTransaction(ctx, func(ctx context.Context) error {
		// 1. Create order
		o, err := s.OrderRepo.Create(ctx, ord)
		if err != nil {
			if errors.Is(err, repository.ErrOrderAlreadyExists) {
				return ErrOrderAlreadyExists
			}
			return err
		}
		created = o

		// 2. Insert items
		items, err := s.ItemsRepo.InsertItems(ctx, o.ID, ord.Items)
		if err != nil {
			return err
		}
		created.Items = items

		// 3. Create outbox event
		ev := entity.OutboxEvent{
			AggregateType: "order",
			AggregateID:   created.ID,
			EventType:     "order.created",
			Payload:       map[string]any{"orderId": created.ID, "userId": created.CustomerID, "totalPrice": created.TotalAmount},
			Status:        entity.OutboxStatus{ID: 1, Name: "pending"},
			CreatedAt:     time.Now(),
		}
		if err := s.OutboxRepo.Create(ctx, ev); err != nil {
			log.Warnf("OrderService.CreateOrder: failed to create outbox: %v", err)
		}

		return nil
	})
	if err != nil {
		log.Errorf("OrderService.CreateOrder: failed: %v", err)
		return entity.Order{}, err
	}

	// 4. Sync Redis
	if err := s.CacheRepo.Save(ctx, &created); err != nil {
		log.Warnf("OrderService.CreateOrder: failed to cache order %s: %v", created.ID, err)
	}
	if err := s.CacheRepo.AddToActive(ctx, &created); err != nil {
		log.Warnf("OrderService.CreateOrder: failed to add to active: %v", err)
	}
	if err := s.CacheRepo.AddUserActive(ctx, created.CustomerID, created.ID); err != nil {
		log.Warnf("OrderService.CreateOrder: failed to add user active: %v", err)
	}

	log.Infof("OrderService.CreateOrder: order %s created", created.ID)
	return created, nil
}

func (s *Service) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status entity.OrderStatus) (entity.Order, error) {
	log.Infof("OrderService.UpdateOrderStatus: order %s -> %s", orderID, status.Name)
	now := time.Now()

	var ord *entity.Order

	err := s.TxManager.WithinTransaction(ctx, func(ctx context.Context) error {
		// Get order from Postgres
		o, err := s.OrderRepo.GetOrderByID(ctx, orderID)
		if err != nil {
			return err
		}
		ord = &o

		// Update in Postgres
		return s.OrderRepo.UpdateOrderStatus(ctx, orderID, status.Name, now)
	})
	if err != nil {
		log.Errorf("OrderService.UpdateOrderStatus: failed: %v", err)
		return entity.Order{}, err
	}

	// Sync redis
	if ord != nil {
		ord.Status = status
		if err := s.CacheRepo.Save(ctx, ord); err != nil {
			log.Warnf("OrderService.UpdateOrderStatus: failed to update cache: %v", err)
		}
		_ = s.CacheRepo.AddToStatus(ctx, string(status.Name), orderID)
		_ = s.CacheRepo.RemoveFromStatus(ctx, string(ord.Status.Name), orderID)
	}

	ord.Status = status

	log.Infof("OrderService.UpdateOrderStatus: order %s updated to %s", orderID, status.Name)
	return *ord, nil
}

func (s *Service) MarkOrderReady(ctx context.Context, orderID uuid.UUID) (entity.Order, error) {
	log.Infof("OrderService.MarkOrderReady: order %s", orderID)
	now := time.Now()
	newStatus := entity.StatusPrepeared

	var ord *entity.Order

	err := s.TxManager.WithinTransaction(ctx, func(ctx context.Context) error {
		// Get order from Postgres
		o, err := s.OrderRepo.GetOrderByID(ctx, orderID)
		if err != nil {
			return err
		}
		ord = &o

		// Update in Postgres
		if err := s.OrderRepo.UpdateOrderStatus(ctx, orderID, newStatus, now); err != nil {
			return err
		}

		// Create outbox event
		ev := entity.OutboxEvent{
			AggregateType: "order",
			AggregateID:   orderID,
			EventType:     "prepeared",
			Payload:       map[string]any{"orderId": orderID},
			Status:        entity.OutboxStatus{Name: entity.OutboxStatusPending},
			CreatedAt:     now,
		}
		if err := s.OutboxRepo.Create(ctx, ev); err != nil {
			log.Warnf("OrderService.MarkOrderReady: failed to create outbox: %v", err)
		}

		return nil
	})
	if err != nil {
		log.Errorf("OrderService.MarkOrderReady: failed: %v", err)
		return entity.Order{}, err
	}

	// Sync redis
	if ord != nil {
		if err := s.CacheRepo.Save(ctx, ord); err != nil {
			log.Warnf("OrderService.MarkOrderReady: failed to update cache: %v", err)
		}
		_ = s.CacheRepo.AddToStatus(ctx, string(newStatus), orderID)
		_ = s.CacheRepo.RemoveFromStatus(ctx, string(ord.Status.Name), orderID)
	}

	ord.Status = entity.OrderStatus{Name: newStatus}

	log.Infof("OrderService.MarkOrderReady: order %s updated to ready", orderID)
	return *ord, nil
}

func (s *Service) MarkOrderPaid(ctx context.Context, orderID, paymentID uuid.UUID) (entity.Order, error) {
	log.Infof("OrderService.MarkOrderPaid: order %s", orderID)
	now := time.Now()

	var ord *entity.Order

	err := s.TxManager.WithinTransaction(ctx, func(ctx context.Context) error {
		// Update order payment and set status
		if err := s.OrderRepo.UpdateOrderPayment(ctx, orderID, paymentID, now); err != nil {
			return err
		}

		// Get updated order with items from Postgres
		o, err := s.OrderRepo.GetOrderByID(ctx, orderID)
		if err != nil {
			return err
		}
		ord = &o

		// Create outbox event
		ev := entity.OutboxEvent{
			AggregateType: "order",
			AggregateID:   orderID,
			EventType:     "paid",
			Payload:       map[string]any{"orderId": orderID, "paymentId": paymentID},
			Status:        entity.OutboxStatus{Name: entity.OutboxStatusPending},
			CreatedAt:     now,
		}
		if err := s.OutboxRepo.Create(ctx, ev); err != nil {
			log.Warnf("OrderService.MarkOrderPaid: failed to create outbox: %v", err)
		}

		return nil
	})
	if err != nil {
		log.Errorf("OrderService.MarkOrderPaid: failed: %v", err)
		return entity.Order{}, err
	}

	// Sync redis
	if ord != nil {
		if err := s.CacheRepo.Save(ctx, ord); err != nil {
			log.Warnf("OrderService.MarkOrderPaid: failed to update cache: %v", err)
		}
		_ = s.CacheRepo.AddToStatus(ctx, string(ord.Status.Name), orderID)
		_ = s.CacheRepo.RemoveFromStatus(ctx, string(entity.StatusCreated), orderID)
	}

	log.Infof("OrderService.MarkOrderPaid: order %s updated", orderID)
	return entity.Order{}, nil
}

func (s *Service) MarkOrderDelivering(ctx context.Context, orderID, deliveryID uuid.UUID) (entity.Order, error) {
	log.Infof("OrderService.MarkOrderDelivering: order %s", orderID)
	now := time.Now()

	var ord *entity.Order

	err := s.TxManager.WithinTransaction(ctx, func(ctx context.Context) error {
		// Update order delivery and set status
		if err := s.OrderRepo.UpdateOrderDelivery(ctx, orderID, deliveryID, now); err != nil {
			return err
		}

		// Get updated order with items from Postgres
		o, err := s.OrderRepo.GetOrderByID(ctx, orderID)
		if err != nil {
			return err
		}
		ord = &o

		// Create outbox event
		ev := entity.OutboxEvent{
			AggregateType: "order",
			AggregateID:   orderID,
			EventType:     "delivering",
			Payload:       map[string]any{"orderId": orderID},
			Status:        entity.OutboxStatus{Name: entity.OutboxStatusPending},
			CreatedAt:     now,
		}
		if err := s.OutboxRepo.Create(ctx, ev); err != nil {
			log.Warnf("OrderService.MarkOrderDelivering: failed to create outbox: %v", err)
		}

		return nil
	})
	if err != nil {
		log.Errorf("OrderService.MarkOrderDelivering: failed: %v", err)
		return entity.Order{}, err
	}

	// Sync redis
	if ord != nil {
		if err := s.CacheRepo.Save(ctx, ord); err != nil {
			log.Warnf("OrderService.MarkOrderPaid: failed to update cache: %v", err)
		}
		_ = s.CacheRepo.AddToStatus(ctx, string(ord.Status.Name), orderID)
		_ = s.CacheRepo.RemoveFromStatus(ctx, string(entity.StatusPrepearing), orderID)
	}

	log.Infof("OrderService.MarkOrderPaid: order %s updated", orderID)
	return entity.Order{}, nil
}

func (s *Service) MarkOrderCompleted(ctx context.Context, orderID uuid.UUID) (entity.Order, error) {
	log.Infof("OrderService.MarkOrderCompleted: order %s", orderID)
	now := time.Now()
	newStatus := entity.StatusCompleted

	var ord *entity.Order

	err := s.TxManager.WithinTransaction(ctx, func(ctx context.Context) error {
		// Update order status
		if err := s.OrderRepo.UpdateOrderStatus(ctx, orderID, newStatus, now); err != nil {
			return err
		}

		// Get updated order with items from Postgres
		o, err := s.OrderRepo.GetOrderByID(ctx, orderID)
		if err != nil {
			return err
		}
		ord = &o

		// Create outbox event
		ev := entity.OutboxEvent{
			AggregateType: "order",
			AggregateID:   orderID,
			EventType:     "completed",
			Payload:       map[string]any{"orderId": orderID},
			Status:        entity.OutboxStatus{Name: entity.OutboxStatusPending},
			CreatedAt:     now,
		}
		if err := s.OutboxRepo.Create(ctx, ev); err != nil {
			log.Warnf("OrderService.MarkOrderCompleted: failed to create outbox: %v", err)
		}

		return nil
	})
	if err != nil {
		log.Errorf("OrderService.MarkOrderCompleted: failed: %v", err)
		return entity.Order{}, err
	}

	// Sync redis
	if ord != nil {
		if err := s.CacheRepo.Save(ctx, ord); err != nil {
			log.Warnf("OrderService.MarkOrderCompleted: failed to update cache: %v", err)
		}
		_ = s.CacheRepo.AddToStatus(ctx, string(ord.Status.Name), orderID)
		_ = s.CacheRepo.RemoveFromStatus(ctx, string(entity.StatusDelivering), orderID)
		_ = s.CacheRepo.RemoveUserActive(ctx, ord.CustomerID, ord.ID)
	}

	log.Infof("OrderService.MarkOrderCompleted: order %s updated", orderID)
	return entity.Order{}, nil
}

func (s *Service) GetOrderByID(ctx context.Context, orderID uuid.UUID) (entity.Order, error) {
	// Attempt to get from cache
	ord, err := s.CacheRepo.GetByID(ctx, orderID)
	if err == nil && ord != nil {
		log.Debugf("OrderService.GetOrderByID: hit cache %s", orderID)
		return *ord, nil
	}

	// Fallback to Postgres
	ordFull, err := s.OrderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		log.Infof("OrderService.GetOrderByID: error: %v", err)
		return entity.Order{}, err
	}

	// Caching
	_ = s.CacheRepo.Save(ctx, &ordFull)
	return ordFull, nil
}

func (s *Service) GetOrdersByUser(ctx context.Context, userID uuid.UUID) ([]entity.Order, error) {
	log.Infof("OrderService.GetOrdersByUser: userID = %v", userID)

	var orders []entity.Order

	ids, err := s.CacheRepo.GetUserActiveOrders(ctx, userID)

	if err == nil && len(ids) > 0 {
		activeOrders := lo.FilterMap(ids, func(idStr string, _ int) (entity.Order, bool) {
			id, err := uuid.Parse(idStr)
			if err != nil {
				return entity.Order{}, false
			}
			ord, err := s.CacheRepo.GetByID(ctx, id)
			if err != nil || ord == nil {
				return entity.Order{}, false
			}
			return *ord, true
		})

		orders = append(orders, activeOrders...)
	}

	inactiveOrders, err := s.OrderRepo.GetOrdersByUserID(ctx, userID)
	if err != nil {
		log.Errorf("OrderService.GetOrdersByUser: error: %v", err)
		return nil, err
	}

	orders = append(orders, inactiveOrders...)

	return orders, nil
}

func (s *Service) GetActiveOrdersByUser(ctx context.Context, userID uuid.UUID) ([]entity.Order, error) {
	log.Infof("OrderService.GetActiveOrdersByUser: userID = %v", userID)

	ids, err := s.CacheRepo.GetUserActiveOrders(ctx, userID)

	if err == nil && len(ids) > 0 {
		activeOrders := lo.FilterMap(ids, func(idStr string, _ int) (entity.Order, bool) {
			id, err := uuid.Parse(idStr)
			if err != nil {
				return entity.Order{}, false
			}
			ord, err := s.CacheRepo.GetByID(ctx, id)
			if err != nil || ord == nil {
				return entity.Order{}, false
			}
			return *ord, true
		})
		if len(activeOrders) > 0 {
			return activeOrders, nil
		}
	}

	return nil, ErrNoActiveOrders
}
