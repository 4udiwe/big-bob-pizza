package payment

import (
	"context"
	"errors"
	"time"

	"github.com/4udiwe/big-bob-pizza/payment-service/internal/entity"
	order_cache_repository "github.com/4udiwe/big-bob-pizza/payment-service/internal/repository/order_cache"
	payment_repository "github.com/4udiwe/big-bob-pizza/payment-service/internal/repository/payment"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func (s *Service) ProcessPayment(ctx context.Context, orderID uuid.UUID, amount float64) (entity.Payment, error) {
	log.Infof("PaymentService.ProcessPayment: orderID=%s amount=%.2f", orderID, amount)

	// 1. Проверяем, что заказ существует и доступен для оплаты
	orderInfo, err := s.OrderCacheRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		if errors.Is(err, order_cache_repository.ErrOrderNotFound) {
			log.Warnf("PaymentService.ProcessPayment: order not found orderID=%s", orderID)
			return entity.Payment{}, ErrOrderNotFound
		}

		log.Errorf("PaymentService.ProcessPayment: failed to get order from cache orderID=%s err=%v", orderID, err)
		return entity.Payment{}, err
	}

	// 2. Проверяем сумму оплаты
	if amount != orderInfo.TotalPrice {
		log.Warnf("PaymentService.ProcessPayment: amount mismatch orderID=%s expected=%.2f got=%.2f",
			orderID, orderInfo.TotalPrice, amount)
		return entity.Payment{}, ErrInvalidAmount
	}

	// 3. Проверяем, не оплачен ли уже заказ
	existingPayment, err := s.PaymentRepo.GetByOrderID(ctx, orderID)
	if err == nil && existingPayment.Status.Name == entity.PaymentStatusCompleted {
		log.Warnf("PaymentService.ProcessPayment: order already paid orderID=%s", orderID)
		return entity.Payment{}, ErrOrderAlreadyPaid
	}

	var payment entity.Payment

	err = s.TxManager.WithinTransaction(ctx, func(ctx context.Context) error {
		// 4. Создаем платеж
		payment = entity.Payment{
			OrderID:  orderID,
			Amount:   amount,
			Currency: "RUB",
			Status:   entity.PaymentStatus{Name: entity.PaymentStatusPending},
		}

		created, err := s.PaymentRepo.Create(ctx, payment)
		if err != nil {
			return err
		}
		payment = created

		// 5. Симулируем обработку платежа (в реальности здесь был бы вызов платежного шлюза)
		// Для демонстрации: 90% успешных платежей, 10% неудачных
		success := time.Now().Unix()%10 != 0 // Простая симуляция

		if success {
			// Успешная оплата
			payment.Status = entity.PaymentStatus{Name: entity.PaymentStatusCompleted}
			if err := s.PaymentRepo.UpdateStatus(ctx, payment.ID, payment.Status, nil); err != nil {
				return err
			}

			// Создаем событие payment.success
			ev := entity.OutboxEvent{
				AggregateType: "payment",
				AggregateID:   payment.ID,
				EventType:     "payment.success",
				Payload: map[string]any{
					"paymentId": payment.ID,
					"orderId":   orderID,
					"amount":    amount,
				},
				Status:    entity.OutboxStatus{Name: entity.OutboxStatusPending},
				CreatedAt: time.Now(),
			}
			if err := s.OutboxRepo.Create(ctx, ev); err != nil {
				log.Warnf("PaymentService.ProcessPayment: failed to create outbox (success): %v", err)
			}
		} else {
			// Неудачная оплата
			reason := "insufficient funds"
			payment.Status = entity.PaymentStatus{Name: entity.PaymentStatusFailed}
			if err := s.PaymentRepo.UpdateStatus(ctx, payment.ID, payment.Status, &reason); err != nil {
				return err
			}

			// Создаем событие payment.failed
			ev := entity.OutboxEvent{
				AggregateType: "payment",
				AggregateID:   payment.ID,
				EventType:     "payment.failed",
				Payload: map[string]any{
					"paymentId": payment.ID,
					"orderId":   orderID,
					"reason":    reason,
				},
				Status:    entity.OutboxStatus{Name: entity.OutboxStatusPending},
				CreatedAt: time.Now(),
			}
			if err := s.OutboxRepo.Create(ctx, ev); err != nil {
				log.Warnf("PaymentService.ProcessPayment: failed to create outbox (failed): %v", err)
			}
		}

		return nil
	})

	if err != nil {
		log.Errorf("PaymentService.ProcessPayment: failed: %v", err)
		return entity.Payment{}, err
	}

	// 6. Удаляем заказ из кэша после обработки платежа
	_ = s.OrderCacheRepo.Delete(ctx, orderID)

	log.Infof("PaymentService.ProcessPayment: payment processed paymentID=%s status=%s", payment.ID, payment.Status.Name)
	return payment, nil
}

func (s *Service) GetPaymentByID(ctx context.Context, paymentID uuid.UUID) (entity.Payment, error) {
	payment, err := s.PaymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		if errors.Is(err, payment_repository.ErrPaymentNotFound) {
			return entity.Payment{}, ErrPaymentNotFound
		}
		return entity.Payment{}, err
	}
	return payment, nil
}

func (s *Service) GetPaymentByOrderID(ctx context.Context, orderID uuid.UUID) (entity.Payment, error) {
	payment, err := s.PaymentRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		if errors.Is(err, payment_repository.ErrPaymentNotFound) {
			return entity.Payment{}, ErrPaymentNotFound
		}
		return entity.Payment{}, err
	}
	return payment, nil
}

func (s *Service) GetAllPayments(ctx context.Context, limit, offset int, status *entity.PaymentStatusName, userID *uuid.UUID) ([]entity.PaymentWithUser, int, error) {
	log.Infof("PaymentService.GetAllPayments: limit=%d offset=%d", limit, offset)
	payments, total, err := s.PaymentRepo.GetAllPayments(ctx, limit, offset, status, userID)
	if err != nil {
		log.Errorf("PaymentService.GetAllPayments: error: %v", err)
		return nil, 0, err
	}
	return payments, total, nil
}

func (s *Service) GetPaymentsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.PaymentWithUser, int, error) {
	log.Infof("PaymentService.GetPaymentsByUserID: userID=%s limit=%d offset=%d", userID, limit, offset)
	payments, total, err := s.PaymentRepo.GetPaymentsByUserID(ctx, userID, limit, offset)
	if err != nil {
		log.Errorf("PaymentService.GetPaymentsByUserID: error: %v", err)
		return nil, 0, err
	}
	return payments, total, nil
}
