package outbox

import (
	"time"

	"github.com/google/uuid"
)

type OutboxEvent struct {
	ID            uuid.UUID
	AggregateType string
	AggregateID   uuid.UUID
	EventType     string
	Payload       map[string]any
	Status        Status
	CreatedAt     time.Time
	ProcessedAt   *time.Time
}
