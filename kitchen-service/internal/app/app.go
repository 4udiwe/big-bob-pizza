package app

import (
	"os"

	"github.com/4udiwe/avito-pvz/pkg/postgres"
	"github.com/labstack/echo"
)

type App struct {
	cfg       *config.Config
	interrupt <-chan os.Signal

	// DB
	postgres *postgres.Postgres
	redis    *redis.Redis

	// Echo
	echoHandler *echo.Echo

	// Repositories
	cacheRepo  *cache_repository.CacheOrderRepository
	orderRepo  *order_repository.Repository
	itemRepo   *item_repository.Repository
	outboxRepo *outbox_repository.Repository

	// Services
	orderService *order.Service

	// Handlers
	postOrderHandler handler.Handler

	// Consumer
	deliveryConsumer *consumer_delivery.Consumer
	kitchenConsumer  *consumer_kitchen.Consumer
	paymentConsumer  *consumer_payment.Consumer

	// Outbox
	OutboxWorker *outbox.Worker
}
