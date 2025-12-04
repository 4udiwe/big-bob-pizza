package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

// KafkaConsumer — тонкая обёртка над kafka-go Reader.
// Позволяет подписаться на один топик и обрабатывать сообщения коллбеком.
type KafkaConsumer struct {
	// brokers — список адресов Kafka‑брокеров.
	brokers []string
	// reader — внутренний kafka-go Reader, создаётся при подписке.
	reader *kafka.Reader
}

// NewConsumer создаёт экземпляр KafkaConsumer c переданными брокерами.
// Подписка на конкретный топик выполняется методом Subscribe.
func NewConsumer(brokers []string) *KafkaConsumer {
	return &KafkaConsumer{
		brokers: brokers,
	}
}

// Subscribe настраивает kafka-go Reader и запускает бесконечный цикл чтения сообщений в отдельной горутине.
//
// Параметры:
//   - ctx — общий контекст сервиса; по его отмене чтение сообщений останавливается;
//   - topic — Kafka‑топик для чтения;
//   - groupID — consumer group id;
//   - handler — пользовательский обработчик сообщения (key, value).
//
// Поведение:
//   - при ошибке FetchMessage и живом контексте — лог, небольшая пауза и повтор;
//   - при ошибке handler — сообщение НЕ коммитится (можно реализовать DLQ отдельно);
//   - при успехе handler — сообщение коммитится.
func (c *KafkaConsumer) Subscribe(
	ctx context.Context,
	topic string,
	groupID string,
	handler func(context.Context, []byte, []byte) error,
) error {
	c.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  c.brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3,
		MaxBytes: 10e6,
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

			// Обработка сообщения пользовательским хендлером.
			if err := handler(ctx, m.Key, m.Value); err != nil {
				logrus.Errorf("KafkaConsumer handler error: %v", err)
				// Ошибка обработки — не коммитим offset, сообщение может быть прочитано снова.
				continue
			}

			// Успешная обработка — коммитим offset.
			if err := c.reader.CommitMessages(ctx, m); err != nil {
				logrus.Errorf("KafkaConsumer commit error: %v", err)
			}
		}
	}()

	return nil
}

// Close вручную закрывает внутренний reader, если он был создан.
// Обычно достаточно отмены контекста, но иногда полезно вызвать явное закрытие.
func (c *KafkaConsumer) Close() error {
	if c.reader != nil {
		return c.reader.Close()
	}
	return nil
}
