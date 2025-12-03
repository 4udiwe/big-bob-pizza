package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
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
			BatchTimeout: 10 * time.Millisecond,
		},
	}
}

func (p *KafkaPublisher) Publish(ctx context.Context, topic string, eventType string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	envelope := Envelope{
		EventID:    uuid.New(),
		EventType:  eventType,
		OccurredAt: time.Now().UTC(),
		Data:       json.RawMessage(data),
	}

	raw, err := json.Marshal(envelope)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(eventType),
		Value: raw,
	}

	return p.writer.WriteMessages(ctx, msg)
}
