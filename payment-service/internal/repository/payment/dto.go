package payment_repository

import (
	"time"

	"github.com/4udiwe/big-bob-pizza/payment-service/internal/entity"
	"github.com/google/uuid"
)

type RowPayment struct {
	ID            uuid.UUID  `db:"id"`
	OrderID       uuid.UUID  `db:"order_id"`
	Amount        float64    `db:"amount"`
	Currency      string     `db:"currency"`
	StatusID      int        `db:"status_id"`
	StatusName    string     `db:"status_name"`
	FailureReason *string    `db:"failure_reason"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
}

func (r RowPayment) ToEntity() entity.Payment {
	return entity.Payment{
		ID:            r.ID,
		OrderID:       r.OrderID,
		Amount:        r.Amount,
		Currency:      r.Currency,
		Status:        entity.PaymentStatus{ID: r.StatusID, Name: entity.PaymentStatusName(r.StatusName)},
		FailureReason: r.FailureReason,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}
}

