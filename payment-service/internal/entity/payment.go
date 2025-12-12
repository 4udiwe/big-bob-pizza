package entity

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatusName string

const (
	PaymentStatusPending   PaymentStatusName = "pending"
	PaymentStatusCompleted PaymentStatusName = "completed"
	PaymentStatusFailed    PaymentStatusName = "failed"
)

type PaymentStatus struct {
	ID   int
	Name PaymentStatusName
}

type Payment struct {
	ID            uuid.UUID
	OrderID       uuid.UUID
	Amount        float64
	Currency      string
	Status        PaymentStatus
	FailureReason *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type PaymentWithUser struct {
	Payment Payment
	UserID  *uuid.UUID
}
