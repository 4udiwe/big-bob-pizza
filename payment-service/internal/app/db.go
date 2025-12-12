package app

import (
	order_cache_repository "github.com/4udiwe/big-bob-pizza/payment-service/internal/repository/order_cache"
	outbox_repository "github.com/4udiwe/big-bob-pizza/payment-service/internal/repository/outbox"
	payment_repository "github.com/4udiwe/big-bob-pizza/payment-service/internal/repository/payment"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/postgres"
)

func (app *App) Postgres() *postgres.Postgres {
	return app.postgres
}

func (app *App) PaymentRepo() *payment_repository.Repository {
	if app.paymentRepo != nil {
		return app.paymentRepo
	}
	app.paymentRepo = payment_repository.New(app.Postgres())
	return app.paymentRepo
}

func (app *App) OrderCacheRepo() *order_cache_repository.Repository {
	if app.orderCacheRepo != nil {
		return app.orderCacheRepo
	}
	app.orderCacheRepo = order_cache_repository.New(app.Postgres())
	return app.orderCacheRepo
}

func (app *App) OutboxRepo() *outbox_repository.Repository {
	if app.outboxRepo != nil {
		return app.outboxRepo
	}
	app.outboxRepo = outbox_repository.New(app.Postgres())
	return app.outboxRepo
}

