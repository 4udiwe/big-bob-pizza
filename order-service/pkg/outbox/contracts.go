package outbox

import (
	"context"

	"github.com/google/uuid"
)

// Repository описывает контракт хранилища для outbox‑записей.
// Конкретная реализация обычно использует таблицу в PostgreSQL.
//
// Типичный жизненный цикл записи:
//   1. Сервис внутри бизнес‑транзакции добавляет запись со статусом "pending".
//   2. Worker периодически вызывает FetchPending и забирает партию таких записей.
//   3. После успешной отправки в брокер вызывается MarkProcessed.
//   4. При ошибке отправки вызывается MarkFailed, а затем периодически RequeueFailed.
type Repository interface {
	// FetchPending возвращает неотправленные события (pending) ограниченным батчем.
	FetchPending(ctx context.Context, limit int) ([]Event, error)
	// MarkProcessed помечает список событий как успешно обработанные.
	MarkProcessed(ctx context.Context, ids []uuid.UUID) error
	// MarkFailed помечает событие как неудавшееся, сохраняя текст ошибки.
	MarkFailed(ctx context.Context, id uuid.UUID, errorText string) error
	// RequeueFailed переводит часть "упавших" событий обратно в pending для повторной отправки.
	RequeueFailed(ctx context.Context, limit int) ([]Event, error)
}

// Publisher описывает транспорт для отправки событий наружу (Kafka, NATS и т.п.).
// В этом проекте реализацией является KafkaPublisher.
type Publisher interface {
	// Publish отправляет произвольный payload в указанный topic с заданным типом события.
	Publish(ctx context.Context, topic string, eventType string, payload any) error
}
