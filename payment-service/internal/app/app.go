package app

import (
	"context"
	"os"

	"github.com/4udiwe/big-bob-pizza/order-service/pkg/kafka"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/outbox"
	"github.com/4udiwe/big-bob-pizza/payment-service/config"
	consumer_order "github.com/4udiwe/big-bob-pizza/payment-service/internal/consumer/order"
	"github.com/4udiwe/big-bob-pizza/payment-service/internal/database"
	"github.com/4udiwe/big-bob-pizza/payment-service/internal/handler"
	order_cache_repository "github.com/4udiwe/big-bob-pizza/payment-service/internal/repository/order_cache"
	outbox_repository "github.com/4udiwe/big-bob-pizza/payment-service/internal/repository/outbox"
	payment_repository "github.com/4udiwe/big-bob-pizza/payment-service/internal/repository/payment"
	"github.com/4udiwe/big-bob-pizza/payment-service/internal/service/payment"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/httpserver"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type App struct {
	cfg       *config.Config
	interrupt <-chan os.Signal

	// DB
	postgres *postgres.Postgres

	// Echo
	echoHandler *echo.Echo

	// Repositories
	paymentRepo    *payment_repository.Repository
	orderCacheRepo *order_cache_repository.Repository
	outboxRepo     *outbox_repository.Repository

	// Services
	paymentService *payment.Service

	// Handlers
	postPaymentHandler handler.Handler

	// Consumer
	orderConsumer *consumer_order.Consumer

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

	// Общий контекст для всех фоновых задач
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Consumer для order.events
	orderKafkaConsumer := kafka.NewConsumer(app.cfg.Kafka.Brokers)

	app.orderConsumer = consumer_order.New(
		app.OrderCacheRepo(),
		orderKafkaConsumer,
		app.cfg.Kafka.Topics.OrderEvents,
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

	// Run consumers and publisher
	app.orderConsumer.Run(ctx)
	app.OutboxWorker.Run(ctx)

	select {
	case s := <-app.interrupt:
		log.Infof("app - Start - signal: %v", s)
	case err := <-httpServer.Notify():
		log.Errorf("app - Start - server error: %v", err)
	}

	// Останавливаем HTTP‑сервер после получения сигнала/ошибки.
	if err := httpServer.Shutdown(); err != nil {
		log.Errorf("HTTP server shutdown error: %v", err)
	}

	log.Info("Shutting down...")
}
