package app

import (
	order_event_repository "github.com/4udiwe/big-bob-pizza/analytics-service/internal/repository/order_event"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/postgres"
)

func (app *App) Postgres() *postgres.Postgres {
	return app.postgres
}

func (app *App) OrderEventRepo() *order_event_repository.Repository {
	if app.orderEventRepo != nil {
		return app.orderEventRepo
	}
	app.orderEventRepo = order_event_repository.New(app.Postgres())
	return app.orderEventRepo
}

