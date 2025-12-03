package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	brokers []string
	reader  *kafka.Reader
}

func NewConsumer(brokers []string) *KafkaConsumer {
	return &KafkaConsumer{
		brokers: brokers,
	}
}

func (c *KafkaConsumer) Subscribe(
	ctx context.Context,
	topic string,
	groupID string,
	handler func(context.Context, []byte, []byte) error,
) error {

	c.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:     c.brokers,
		GroupID:     groupID,
		Topic:       topic,
		MinBytes:    10e3,
		MaxBytes:    10e6,
	})

	go func() {
		defer c.reader.Close()

		for {
			m, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					logrus.Info("KafkaConsumer: context cancelled, stopping...")
					return
				}
				logrus.Errorf("KafkaConsumer fetch error: %v", err)
				time.Sleep(time.Second)
				continue
			}

			// Handling
			if err := handler(ctx, m.Key, m.Value); err != nil {
				logrus.Errorf("KafkaConsumer handler error: %v", err)
				// no commit if error
				continue
			}

			// commit if success
			if err := c.reader.CommitMessages(ctx, m); err != nil {
				logrus.Errorf("KafkaConsumer commit error: %v", err)
			}
		}
	}()

	return nil
}

func (c *KafkaConsumer) Close() error {
	if c.reader != nil {
		return c.reader.Close()
	}
	return nil
}
