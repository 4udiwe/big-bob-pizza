package app

import (
	"github.com/4udiwe/big-bob-pizza/payment-service/internal/handler"
	"github.com/4udiwe/big-bob-pizza/payment-service/internal/handler/post_payment"
)

func (app *App) PostPaymentHandler() handler.Handler {
	if app.postPaymentHandler != nil {
		return app.postPaymentHandler
	}
	app.postPaymentHandler = post_payment.New(app.PaymentService())
	return app.postPaymentHandler
}

