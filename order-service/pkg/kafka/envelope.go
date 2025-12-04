package kafka

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Envelope — единый формат обёртки для всех событий, публикуемых в Kafka.
// Это позволяет:
//   - иметь общий EventID для трейсинга;
//   - хранить время возникновения события отдельно от времени публикации;
//   - иметь тип события (EventType) и произвольный payload (Data).
type Envelope struct {
	// EventID — идентификатор события (генерируется при публикации).
	EventID uuid.UUID `json:"eventId"`
	// EventType — строковый тип события (например, "OrderCreated").
	EventType string `json:"eventType"`
	// OccurredAt — момент времени, когда событие произошло в доменной модели.
	OccurredAt time.Time `json:"occuredAt"`
	// Data — сырое тело события в виде JSON (конкретный payload доменного события).
	Data json.RawMessage `json:"data"`
}
