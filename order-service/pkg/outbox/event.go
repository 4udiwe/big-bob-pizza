package outbox

import "github.com/google/uuid"

// Event — минимальное представление записи в outbox‑таблице.
//   - ID — идентификатор записи в outbox (обычно UUID из БД);
//   - EventType — тип доменного события (например, "OrderCreated");
//   - Payload — сериализованное тело события (JSON, protobuf и т.п.).
type Event struct {
	ID        uuid.UUID
	EventType string
	Payload   []byte
}
