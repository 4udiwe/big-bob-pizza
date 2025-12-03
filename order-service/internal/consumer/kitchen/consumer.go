package consumer_kitchen

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/4udiwe/big-bob-pizza/order-service/internal/consumer"
	"github.com/4udiwe/big-bob-pizza/order-service/internal/entity"
	"github.com/4udiwe/big-bob-pizza/order-service/internal/service/order"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/kafka"
)

// Обработчик событий для топика кухни
type Consumer struct {
	svc      *order.Service
	consumer *kafka.KafkaConsumer
	topic    string
	groupID  string
}

func New(
	svc *order.Service,
	consumer *kafka.KafkaConsumer,
	topic string,
	groupID string,
) *Consumer {
	return &Consumer{
		svc:      svc,
		consumer: consumer,
		topic:    topic,
		groupID:  groupID,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	logrus.Infof("OrderConsumer: subscribing to topic=%s group=%s", c.topic, c.groupID)

	return c.consumer.Subscribe(ctx, c.topic, c.groupID, func(ctx context.Context, key, value []byte) error {
		event, err := consumer.ParseOrderEvent(value)
		if err != nil {
			logrus.Errorf("OrderConsumer: failed to parse event: %v", err)
			return nil
		}

		switch event.Type {
		case consumer.KitchenAccepted:
			status := entity.OrderStatus{Name: entity.StatusPrepearing}
			_, err = c.svc.UpdateOrderStatus(ctx, event.Payload.OrderID, status)
			if err != nil {
				logrus.Errorf("OrderConsumer: UpdateOrderStatus(prepearing) failed: %v", err)
			}

		case consumer.KitchenReady:
			_, err := c.svc.MarkOrderReady(ctx, event.Payload.OrderID)
			if err != nil {
				logrus.Errorf("OrderConsumer: MarkOrderReady failed: %v", err)
			}

		case consumer.KitchenHandedToCourier:
			_, err = c.svc.MarkOrderDelivering(ctx, event.Payload.OrderID, event.Payload.DeliveryID)
			if err != nil {
				logrus.Errorf("OrderConsumer: MarkOrderDelivering failed: %v", err)
			}

		default:
			logrus.Errorf("OrderConsumer: unknown event type %s", event.Type)
			return nil
		}

		return err
	})
}
