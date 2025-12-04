package app

import (
	cache_repository "github.com/4udiwe/big-bob-pizza/order-service/internal/repository/cache"
	item_repository "github.com/4udiwe/big-bob-pizza/order-service/internal/repository/item"
	order_repository "github.com/4udiwe/big-bob-pizza/order-service/internal/repository/order"
	outbox_repository "github.com/4udiwe/big-bob-pizza/order-service/internal/repository/outbox"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/postgres"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/redis"
)

func (app *App) Postgres() *postgres.Postgres {
	return app.postgres
}

func (app *App) Redis() *redis.Redis {
	return app.redis
}

func (app *App) CacheRepo() *cache_repository.CacheOrderRepository {
	if app.cacheRepo != nil {
		return app.cacheRepo
	}
	app.cacheRepo = cache_repository.NewCacheOrderRepository(app.Redis())
	return app.cacheRepo
}

func (app *App) OrderRepo() *order_repository.Repository {
	if app.orderRepo != nil {
		return app.orderRepo
	}
	app.orderRepo = order_repository.New(app.Postgres())
	return app.orderRepo
}

func (app *App) ItemRepo() *item_repository.Repository {
	if app.itemRepo != nil {
		return app.itemRepo
	}
	app.itemRepo = item_repository.New(app.Postgres())
	return app.itemRepo
}

func (app *App) OutboxRepo() *outbox_repository.Repository {
	if app.outboxRepo != nil {
		return app.outboxRepo
	}
	app.outboxRepo = outbox_repository.New(app.Postgres())
	return app.outboxRepo
}
