package order_event

import (
	"database/sql"
	"time"

	"github.com/4udiwe/big-bob-pizza/analytics-service/internal/entity"
	"github.com/google/uuid"
)

type OrderEventDTO struct {
	ID         uuid.UUID
	EventID    uuid.UUID
	EventType  string
	OrderID    uuid.UUID
	UserID     sql.NullString
	Amount     sql.NullFloat64
	PaymentID  sql.NullString
	Reason     sql.NullString
	OccurredAt time.Time
	CreatedAt  time.Time
}

func (dto *OrderEventDTO) ToEntity() entity.OrderEvent {
	event := entity.OrderEvent{
		ID:         dto.ID,
		EventID:    dto.EventID,
		EventType:  dto.EventType,
		OrderID:    dto.OrderID,
		OccurredAt: dto.OccurredAt,
		CreatedAt:  dto.CreatedAt,
	}

	if dto.UserID.Valid {
		if uid, err := uuid.Parse(dto.UserID.String); err == nil {
			event.UserID = &uid
		}
	}

	if dto.Amount.Valid {
		event.Amount = &dto.Amount.Float64
	}

	if dto.PaymentID.Valid {
		if pid, err := uuid.Parse(dto.PaymentID.String); err == nil {
			event.PaymentID = &pid
		}
	}

	if dto.Reason.Valid {
		event.Reason = &dto.Reason.String
	}

	return event
}

func (dto *OrderEventDTO) FromEntity(event entity.OrderEvent) {
	dto.ID = event.ID
	dto.EventID = event.EventID
	dto.EventType = event.EventType
	dto.OrderID = event.OrderID
	dto.OccurredAt = event.OccurredAt
	dto.CreatedAt = event.CreatedAt

	if event.UserID != nil {
		dto.UserID = sql.NullString{String: event.UserID.String(), Valid: true}
	}

	if event.Amount != nil {
		dto.Amount = sql.NullFloat64{Float64: *event.Amount, Valid: true}
	}

	if event.PaymentID != nil {
		dto.PaymentID = sql.NullString{String: event.PaymentID.String(), Valid: true}
	}

	if event.Reason != nil {
		dto.Reason = sql.NullString{String: *event.Reason, Valid: true}
	}
}

