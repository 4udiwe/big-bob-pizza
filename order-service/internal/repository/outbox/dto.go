package outbox_repository

import (
	"time"

	"github.com/4udiwe/big-bob-pizza/order-service/internal/entity/outbox"
	"github.com/google/uuid"
)

type RowOutbox struct {
	ID            uuid.UUID      `db:"id"`
	AggregateType string         `db:"aggregate_type"`
	AggregateID   uuid.UUID      `db:"aggregate_id"`
	EventType     string         `db:"event_type"`
	Payload       map[string]any `db:"payload"`
	StatusID      int            `db:"status_id"`
	StatusName    string         `db:"status_name"`
	CreatedAt     time.Time      `db:"created_at"`
	ProcessedAt   *time.Time     `db:"processed_at"`
}

func (r RowOutbox) ToEntity() outbox.OutboxEvent {
	return outbox.OutboxEvent{
		ID:            r.ID,
		AggregateType: r.AggregateType,
		AggregateID:   r.AggregateID,
		EventType:     r.EventType,
		Payload:       r.Payload,
		Status:        outbox.Status{ID: r.StatusID, Name: outbox.StatusName(r.StatusName)},
		CreatedAt:     r.CreatedAt,
		ProcessedAt:   r.ProcessedAt,
	}
}
