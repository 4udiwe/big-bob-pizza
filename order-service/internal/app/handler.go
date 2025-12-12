package app

import (
	"github.com/4udiwe/big-bob-pizza/order-service/internal/handler"
	"github.com/4udiwe/big-bob-pizza/order-service/internal/handler/post_order"
)

func (app *App) PostOrderHandler() handler.Handler {
	if app.postOrderHandler != nil {
		return app.postOrderHandler
	}
	app.postOrderHandler = post_order.New(app.OrderService())
	return app.postOrderHandler
}
