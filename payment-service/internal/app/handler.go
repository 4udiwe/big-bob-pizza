package app

import (
	"github.com/4udiwe/big-bob-pizza/payment-service/internal/handler"
	get_payment "github.com/4udiwe/big-bob-pizza/payment-service/internal/handler/get_payment"
	get_payment_by_order "github.com/4udiwe/big-bob-pizza/payment-service/internal/handler/get_payment_by_order"
	get_payments "github.com/4udiwe/big-bob-pizza/payment-service/internal/handler/get_payments"
	"github.com/4udiwe/big-bob-pizza/payment-service/internal/handler/post_payment"
)

func (app *App) PostPaymentHandler() handler.Handler {
	if app.postPaymentHandler != nil {
		return app.postPaymentHandler
	}
	app.postPaymentHandler = post_payment.New(app.PaymentService())
	return app.postPaymentHandler
}

func (app *App) GetPaymentsHandler() handler.Handler {
	return get_payments.New(app.PaymentService())
}

func (app *App) GetPaymentHandler() handler.Handler {
	return get_payment.New(app.PaymentService())
}

func (app *App) GetPaymentByOrderHandler() handler.Handler {
	return get_payment_by_order.New(app.PaymentService())
}

