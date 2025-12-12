package consumer_order

import (
	"context"
	"encoding/json"

	"github.com/4udiwe/big-bob-pizza/order-service/pkg/kafka"
	"github.com/4udiwe/big-bob-pizza/payment-service/internal/entity"
	order_cache "github.com/4udiwe/big-bob-pizza/payment-service/internal/repository/order_cache"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Consumer обрабатывает события из топика order.events
type Consumer struct {
	orderCacheRepo *order_cache.Repository
	consumer       *kafka.KafkaConsumer
	topic          string
	groupID        string
}

func New(
	orderCacheRepo *order_cache.Repository,
	consumer *kafka.KafkaConsumer,
	topic string,
	groupID string,
) *Consumer {
	return &Consumer{
		orderCacheRepo: orderCacheRepo,
		consumer:       consumer,
		topic:          topic,
		groupID:        groupID,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	logrus.Infof("OrderConsumer: subscribing to topic=%s group=%s", c.topic, c.groupID)

	return c.consumer.Subscribe(ctx, c.topic, c.groupID, func(ctx context.Context, key, value []byte) error {
		// Парсим envelope
		var env kafka.Envelope
		if err := json.Unmarshal(value, &env); err != nil {
			logrus.Errorf("OrderConsumer: failed to parse envelope: %v", err)
			return nil
		}

		// Обрабатываем только событие order.created
		if env.EventType != "order.created" {
			return nil
		}

		// Парсим payload события order.created
		var payload struct {
			OrderID    uuid.UUID `json:"orderId"`
			UserID     uuid.UUID `json:"userId"`
			TotalPrice float64   `json:"totalPrice"`
		}

		if env.Data == nil {
			logrus.Errorf("OrderConsumer: empty data for event %s", env.EventType)
			return nil
		}

		if err := json.Unmarshal(env.Data, &payload); err != nil {
			logrus.Errorf("OrderConsumer: failed to parse payload: %v", err)
			return nil
		}

		// Сохраняем информацию о заказе в кэш для последующей оплаты
		orderInfo := entity.OrderInfo{
			OrderID:    payload.OrderID,
			UserID:     payload.UserID,
			TotalPrice: payload.TotalPrice,
			CreatedAt:  env.OccurredAt,
		}

		if err := c.orderCacheRepo.Save(ctx, orderInfo); err != nil {
			logrus.Errorf("OrderConsumer: failed to save order to cache: %v", err)
			return err
		}

		logrus.Infof("OrderConsumer: order cousumed and cached for payment orderID=%s", payload.OrderID)
		return nil
	})
}
