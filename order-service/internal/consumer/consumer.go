package consumer

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/4udiwe/big-bob-pizza/order-service/internal/consumer/events"
	"github.com/4udiwe/big-bob-pizza/order-service/internal/entity"
	"github.com/4udiwe/big-bob-pizza/order-service/internal/service/order"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/kafka"
)

type OrderEvent struct {
	Type       string    `json:"type"`
	OrderID    uuid.UUID `json:"order_id"`
	PaymentID  uuid.UUID `json:"payment_id,omitempty"`
	DeliveryID uuid.UUID `json:"delivery_id,omitempty"`
}

type OrderConsumer struct {
	svc      *order.Service
	consumer *kafka.KafkaConsumer
	topic    string
	groupID  string
}

func NewOrderConsumer(
	svc *order.Service,
	consumer *kafka.KafkaConsumer,
	topic string,
	groupID string,
) *OrderConsumer {
	return &OrderConsumer{
		svc:      svc,
		consumer: consumer,
		topic:    topic,
		groupID:  groupID,
	}
}

func (c *OrderConsumer) Run(ctx context.Context) error {
	logrus.Infof("OrderConsumer: subscribing to topic=%s group=%s", c.topic, c.groupID)

	return c.consumer.Subscribe(ctx, c.topic, c.groupID, func(ctx context.Context, key, value []byte) error {
		event, err := events.ParseOrderEvent(value)
		if err != nil {
			logrus.Errorf("OrderConsumer: failed to parse event: %v", err)
			return nil
		}

		switch event.Type {
		case events.PaymentSuccess:
			_, err = c.svc.MarkOrderPaid(ctx, event.Payload.OrderID, event.Payload.PaymentID)
			if err != nil {
				logrus.Errorf("OrderConsumer: MarkOrderPaid failed: %v", err)
			}

		case events.PaymentFailed:
			status := entity.OrderStatus{Name: entity.StatusCancelled}
			_, err = c.svc.UpdateOrderStatus(ctx, event.Payload.OrderID, status)
			if err != nil {
				logrus.Errorf("OrderConsumer: UpdateOrderStatus(cancelled) failed: %v", err)
			}

		case events.KitchenAccepted:
			status := entity.OrderStatus{Name: entity.StatusPrepearing}
			_, err = c.svc.UpdateOrderStatus(ctx, event.Payload.OrderID, status)
			if err != nil {
				logrus.Errorf("OrderConsumer: UpdateOrderStatus(prepearing) failed: %v", err)
			}

		case events.KitchenReady:
			_, err := c.svc.MarkOrderReady(ctx, event.Payload.OrderID)
			if err != nil {
				logrus.Errorf("OrderConsumer: MarkOrderReady failed: %v", err)
			}

		case events.KitchenHandedToCourier:
			_, err = c.svc.MarkOrderDelivering(ctx, event.Payload.OrderID, event.Payload.DeliveryID)
			if err != nil {
				logrus.Errorf("OrderConsumer: MarkOrderDelivering failed: %v", err)
			}

		case events.DeliveryCompleted:
			_, err = c.svc.MarkOrderCompleted(ctx, event.Payload.OrderID)
			if err != nil {
				logrus.Errorf("OrderConsumer: MarkOrderCompleted failed: %v", err)
			}

		default:
			logrus.Errorf("OrderConsumer: unknown event type %s", event.Type)
			return nil
		}

		return err
	})
}
