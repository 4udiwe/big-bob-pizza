package entity

import (
	"time"

	"github.com/google/uuid"
)

type OutboxStatusName string

const (
	OutboxStatusPending   OutboxStatusName = "pending"
	OutboxStatusFailed    OutboxStatusName = "failed"
	OutboxStatusProcessed OutboxStatusName = "processed"
)

type OutboxStatus struct {
	ID   int
	Name OutboxStatusName
}

type OutboxEvent struct {
	ID            uuid.UUID
	AggregateType string
	AggregateID   uuid.UUID
	EventType     string
	Payload       map[string]any
	Status        OutboxStatus
	CreatedAt     time.Time
	ProcessedAt   *time.Time
}
