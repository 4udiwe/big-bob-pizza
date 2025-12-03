package outbox_repository

import (
	"encoding/json"
	"time"

	"github.com/4udiwe/big-bob-pizza/order-service/internal/entity"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/outbox"
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

func (r RowOutbox) ToEntity() entity.OutboxEvent {
	return entity.OutboxEvent{
		ID:            r.ID,
		AggregateType: r.AggregateType,
		AggregateID:   r.AggregateID,
		EventType:     r.EventType,
		Payload:       r.Payload,
		Status:        entity.OutboxStatus{ID: r.StatusID, Name: entity.OutboxStatusName(r.StatusName)},
		CreatedAt:     r.CreatedAt,
		ProcessedAt:   r.ProcessedAt,
	}
}

func (r RowOutbox) ToEvent() outbox.Event {
	payloadBytes, err := json.Marshal(r.Payload)
	if err != nil {
		payloadBytes = []byte("{}")
	}
	return outbox.Event{
		ID:        r.ID,
		EventType: r.EventType,
		Payload:   payloadBytes,
	}
}
