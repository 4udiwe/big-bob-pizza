package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

// KafkaPublisher — обёртка над kafka-go Writer для публикации событий в формате Envelope.
type KafkaPublisher struct {
	writer *kafka.Writer
}

// NewKafkaPublisher создаёт синхронный Kafka‑паблишер с минимальными настройками,
// используя переданный список брокеров.
func NewKafkaPublisher(brokers []string) *KafkaPublisher {
	return &KafkaPublisher{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Async:        false,
			BatchTimeout: 10 * time.Millisecond,
		},
	}
}

// Publish сериализует произвольный payload в JSON, заворачивает его в Envelope
// и публикует в Kafka в указанный topic.
//
// Ключ сообщения (Key) — это eventType, чтобы события одного типа лежали последовательно в партициях.
func (p *KafkaPublisher) Publish(ctx context.Context, topic string, eventType string, payload any) error {
	var raw json.RawMessage
	switch v := payload.(type) {
	case json.RawMessage:
		raw = v
	case []byte:
		raw = json.RawMessage(v)
	default:
		data, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		raw = json.RawMessage(data)
	}

	envelope := Envelope{
		EventID:    uuid.New(),
		EventType:  eventType,
		OccurredAt: time.Now().UTC(),
		Data:       raw,
	}

	raw, err := json.Marshal(&envelope)
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
