package payment

import "errors"

var (
	ErrOrderNotFound    = errors.New("order not found or expired")
	ErrOrderAlreadyPaid = errors.New("order already paid")
	ErrPaymentNotFound  = errors.New("payment not found")
	ErrInvalidAmount    = errors.New("payment amount does not match order amount")
)

