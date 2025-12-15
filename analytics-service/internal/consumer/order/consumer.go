package consumer_order

import (
	"context"
	"encoding/json"

	"github.com/4udiwe/big-bob-pizza/analytics-service/internal/entity"
	"github.com/4udiwe/big-bob-pizza/analytics-service/internal/service/analytics"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/kafka"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Consumer обрабатывает события из топика order.events
type Consumer struct {
	analyticsService *analytics.Service
	consumer         *kafka.KafkaConsumer
	topic            string
	groupID          string
}

func New(
	analyticsService *analytics.Service,
	consumer *kafka.KafkaConsumer,
	topic string,
	groupID string,
) *Consumer {
	return &Consumer{
		analyticsService: analyticsService,
		consumer:         consumer,
		topic:            topic,
		groupID:          groupID,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	logrus.Infof("OrderAnalyticsConsumer: subscribing to topic=%s group=%s", c.topic, c.groupID)

	return c.consumer.Subscribe(ctx, c.topic, c.groupID, func(ctx context.Context, key, value []byte) error {
		// Парсим envelope
		var env kafka.Envelope
		if err := json.Unmarshal(value, &env); err != nil {
			logrus.Errorf("OrderAnalyticsConsumer: failed to parse envelope: %v", err)
			return nil
		}

		// Обрабатываем только нужные события
		switch env.EventType {
		case "order.created":
			return c.handleOrderCreated(ctx, env)
		case "order.paid":
			return c.handleOrderPaid(ctx, env)
		case "order.cancelled":
			return c.handleOrderCancelled(ctx, env)
		case "order.completed":
			return c.handleOrderCompleted(ctx, env)
		default:
			// Игнорируем другие события
			return nil
		}
	})
}

func (c *Consumer) handleOrderCreated(ctx context.Context, env kafka.Envelope) error {
	var payload struct {
		OrderID    uuid.UUID `json:"orderId"`
		UserID     uuid.UUID `json:"userId"`
		TotalPrice float64   `json:"totalPrice"`
	}

	if env.Data == nil {
		logrus.Errorf("OrderAnalyticsConsumer: empty data for event %s", env.EventType)
		return nil
	}

	if err := json.Unmarshal(env.Data, &payload); err != nil {
		logrus.Errorf("OrderAnalyticsConsumer: failed to parse payload: %v", err)
		return nil
	}

	event := entity.OrderEvent{
		EventID:    env.EventID,
		EventType:  "order.created",
		OrderID:    payload.OrderID,
		UserID:     &payload.UserID,
		Amount:     &payload.TotalPrice,
		OccurredAt: env.OccurredAt,
	}

	if err := c.analyticsService.SaveOrderEvent(ctx, event); err != nil {
		logrus.Errorf("OrderAnalyticsConsumer: failed to save order.created event: %v", err)
		return err
	}

	logrus.Infof("OrderAnalyticsConsumer: processed order.created orderID=%s", payload.OrderID)
	return nil
}

func (c *Consumer) handleOrderPaid(ctx context.Context, env kafka.Envelope) error {
	var payload struct {
		OrderID   uuid.UUID `json:"orderId"`
		PaymentID uuid.UUID `json:"paymentId"`
	}

	if env.Data == nil {
		logrus.Errorf("OrderAnalyticsConsumer: empty data for event %s", env.EventType)
		return nil
	}

	if err := json.Unmarshal(env.Data, &payload); err != nil {
		logrus.Errorf("OrderAnalyticsConsumer: failed to parse payload: %v", err)
		return nil
	}

	event := entity.OrderEvent{
		EventID:    env.EventID,
		EventType:  "order.paid",
		OrderID:    payload.OrderID,
		PaymentID:  &payload.PaymentID,
		OccurredAt: env.OccurredAt,
	}

	if err := c.analyticsService.SaveOrderEvent(ctx, event); err != nil {
		logrus.Errorf("OrderAnalyticsConsumer: failed to save order.paid event: %v", err)
		return err
	}

	logrus.Infof("OrderAnalyticsConsumer: processed order.paid orderID=%s", payload.OrderID)
	return nil
}

func (c *Consumer) handleOrderCancelled(ctx context.Context, env kafka.Envelope) error {
	var payload struct {
		OrderID uuid.UUID `json:"orderId"`
		Reason  string    `json:"reason"`
	}

	if env.Data == nil {
		logrus.Errorf("OrderAnalyticsConsumer: empty data for event %s", env.EventType)
		return nil
	}

	if err := json.Unmarshal(env.Data, &payload); err != nil {
		logrus.Errorf("OrderAnalyticsConsumer: failed to parse payload: %v", err)
		return nil
	}

	event := entity.OrderEvent{
		EventID:    env.EventID,
		EventType:  "order.cancelled",
		OrderID:    payload.OrderID,
		Reason:     &payload.Reason,
		OccurredAt: env.OccurredAt,
	}

	if err := c.analyticsService.SaveOrderEvent(ctx, event); err != nil {
		logrus.Errorf("OrderAnalyticsConsumer: failed to save order.cancelled event: %v", err)
		return err
	}

	logrus.Infof("OrderAnalyticsConsumer: processed order.cancelled orderID=%s", payload.OrderID)
	return nil
}

func (c *Consumer) handleOrderCompleted(ctx context.Context, env kafka.Envelope) error {
	var payload struct {
		OrderID uuid.UUID `json:"orderId"`
	}

	if env.Data == nil {
		logrus.Errorf("OrderAnalyticsConsumer: empty data for event %s", env.EventType)
		return nil
	}

	if err := json.Unmarshal(env.Data, &payload); err != nil {
		logrus.Errorf("OrderAnalyticsConsumer: failed to parse payload: %v", err)
		return nil
	}

	event := entity.OrderEvent{
		EventID:    env.EventID,
		EventType:  "order.completed",
		OrderID:    payload.OrderID,
		OccurredAt: env.OccurredAt,
	}

	if err := c.analyticsService.SaveOrderEvent(ctx, event); err != nil {
		logrus.Errorf("OrderAnalyticsConsumer: failed to save order.completed event: %v", err)
		return err
	}

	logrus.Infof("OrderAnalyticsConsumer: processed order.completed orderID=%s", payload.OrderID)
	return nil
}

