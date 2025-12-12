package app

import (
	"github.com/4udiwe/big-bob-pizza/order-service/internal/handler"
	get_active_orders_by_user "github.com/4udiwe/big-bob-pizza/order-service/internal/handler/get_active_orders_by_user"
	get_all_orders "github.com/4udiwe/big-bob-pizza/order-service/internal/handler/get_all_orders"
	get_order "github.com/4udiwe/big-bob-pizza/order-service/internal/handler/get_order"
	get_orders_by_user "github.com/4udiwe/big-bob-pizza/order-service/internal/handler/get_orders_by_user"
	"github.com/4udiwe/big-bob-pizza/order-service/internal/handler/post_order"
)

func (app *App) PostOrderHandler() handler.Handler {
	if app.postOrderHandler != nil {
		return app.postOrderHandler
	}
	app.postOrderHandler = post_order.New(app.OrderService())
	return app.postOrderHandler
}

func (app *App) GetOrderHandler() handler.Handler {
	return get_order.New(app.OrderService())
}

func (app *App) GetOrdersByUserHandler() handler.Handler {
	return get_orders_by_user.New(app.OrderService())
}

func (app *App) GetActiveOrdersByUserHandler() handler.Handler {
	return get_active_orders_by_user.New(app.OrderService())
}

func (app *App) GetAllOrdersHandler() handler.Handler {
	return get_all_orders.New(app.OrderService())
}
