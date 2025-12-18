package app

import (
	"github.com/4udiwe/big-bob-pizza/payment-service/internal/service/payment"
)

func (app *App) PaymentService() *payment.Service {
	if app.paymentService != nil {
		return app.paymentService
	}
	app.paymentService = payment.NewService(
		app.PaymentRepo(),
		app.OrderCacheRepo(),
		app.OutboxRepo(),
		app.Postgres(),
	)
	return app.paymentService
}

