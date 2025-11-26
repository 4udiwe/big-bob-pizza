package kafka_publicher

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type KafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher(brokers []string) *KafkaPublisher {
	return &KafkaPublisher{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Async:        false,
			BatchTimeout: 10,
		},
	}
}

func (p *KafkaPublisher) Publish(ctx context.Context, topic string, payload []byte) error {
	msg := kafka.Message{
		Topic: topic,
		Value: payload,
	}

	return p.writer.WriteMessages(ctx, msg)
}
