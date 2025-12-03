package app

import (
	"context"
	"os"

	"github.com/4udiwe/big-bob-pizza/order-service/config"
	consumer_delivery "github.com/4udiwe/big-bob-pizza/order-service/internal/consumer/delivery"
	consumer_kitchen "github.com/4udiwe/big-bob-pizza/order-service/internal/consumer/kitchen"
	consumer_payment "github.com/4udiwe/big-bob-pizza/order-service/internal/consumer/payment"
	"github.com/4udiwe/big-bob-pizza/order-service/internal/database"
	"github.com/4udiwe/big-bob-pizza/order-service/internal/handler"
	cache_repository "github.com/4udiwe/big-bob-pizza/order-service/internal/repository/cache"
	item_repository "github.com/4udiwe/big-bob-pizza/order-service/internal/repository/item"
	order_repository "github.com/4udiwe/big-bob-pizza/order-service/internal/repository/order"
	outbox_repository "github.com/4udiwe/big-bob-pizza/order-service/internal/repository/outbox"
	"github.com/4udiwe/big-bob-pizza/order-service/internal/service/order"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/httpserver"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/kafka"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/outbox"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/postgres"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/redis"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
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

func New(configPath string) *App {
	cfg, err := config.New(configPath)
	if err != nil {
		log.Fatalf("app - New - config.New: %v", err)
	}

	initLogger(cfg.Log.Level)

	return &App{
		cfg: cfg,
	}
}

func (app *App) Start() {
	// Postgres
	log.Info("Connecting to PostgreSQL...")

	postgres, err := postgres.New(app.cfg.Postgres.URL, postgres.ConnAttempts(5))

	if err != nil {
		log.Fatalf("app - Start - Postgres failed:%v", err)
	}
	app.postgres = postgres

	defer postgres.Close()

	// Migrations
	if err := database.RunMigrations(context.Background(), app.postgres.Pool); err != nil {
		log.Errorf("app - Start - Migrations failed: %v", err)
	}

	// Redis
	redis, err := redis.New(app.cfg.Redis.Addr, "", 0)
	if err != nil {
		log.Fatalf("app - Start - Redis failed:%v", err)
	}
	app.redis = redis

	defer redis.Close()

	// Consumers
	paymentKafkaConsumer := kafka.NewConsumer(app.cfg.Kafka.Brokers)
	kitchenKafkaConsumer := kafka.NewConsumer(app.cfg.Kafka.Brokers)
	deliveryKafkaConsumer := kafka.NewConsumer(app.cfg.Kafka.Brokers)

	app.paymentConsumer = consumer_payment.New(
		app.orderService,
		paymentKafkaConsumer,
		app.cfg.Kafka.Topics.PaymentEvents,
		app.cfg.Kafka.Consumer.GroupID,
	)

	app.kitchenConsumer = consumer_kitchen.New(
		app.orderService,
		kitchenKafkaConsumer,
		app.cfg.Kafka.Topics.KitchenEvents,
		app.cfg.Kafka.Consumer.GroupID,
	)

	app.deliveryConsumer = consumer_delivery.New(
		app.orderService,
		deliveryKafkaConsumer,
		app.cfg.Kafka.Topics.DeliveryEvents,
		app.cfg.Kafka.Consumer.GroupID,
	)

	// Outbox publisher
	kafkaPublisher := kafka.NewKafkaPublisher(app.cfg.Kafka.Brokers)

	app.OutboxWorker = outbox.NewWorker(
		app.OutboxRepo(),
		kafkaPublisher,
		app.cfg.Outbox.Topic,
		app.cfg.Outbox.BatchLimit,
		app.cfg.Outbox.RequeBatchLimit,
		app.cfg.Outbox.Interval,
		app.cfg.Outbox.RequeInterval,
	)

	// App server
	log.Info("Starting app server...")
	httpServer := httpserver.New(app.EchoHandler(), httpserver.Port(app.cfg.HTTP.Port))
	httpServer.Start()
	log.Debugf("Server port: %s", app.cfg.HTTP.Port)

	defer func() {
		if err := httpServer.Shutdown(); err != nil {
			log.Errorf("HTTP server shutdown error: %v", err)
		}
	}()

	// Run consumers and publisher
	ctx := context.Background()

	app.paymentConsumer.Run(ctx)
	app.kitchenConsumer.Run(ctx)
	app.deliveryConsumer.Run(ctx)

	app.OutboxWorker.Run(ctx)

	select {
	case s := <-app.interrupt:
		log.Infof("app - Start - signal: %v", s)
	case err := <-httpServer.Notify():
		log.Errorf("app - Start - server error: %v", err)
	}

	log.Info("Shutting down...")
}
