package consumer_order

import (
	"context"
	"encoding/json"

	"github.com/4udiwe/big-bob-pizza/kitchen-service/internal/entity"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/kafka"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Consumer обрабатывает события из топика order.events
type Consumer struct {
	service  *Service
	consumer *kafka.KafkaConsumer
	topic    string
	groupID  string
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
		case "order.paid":
			return c.handleOrderPaid(ctx, env)
		default:
			// Игнорируем другие события
			return nil
		}
	})
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
