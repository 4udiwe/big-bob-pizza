package app

import "github.com/4udiwe/big-bob-pizza/order-service/internal/service/order"

func (app *App) OrderService() *order.Service {
	if app.orderService != nil {
		return app.orderService
	}
	app.orderService = order.NewService(
		app.OrderRepo(),
		app.ItemRepo(),
		app.OutboxRepo(),
		app.CacheRepo(),
		app.Postgres(),
	)
	return app.orderService
}
