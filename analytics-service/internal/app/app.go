package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/4udiwe/big-bob-pizza/analytics-service/config"
	consumer_order "github.com/4udiwe/big-bob-pizza/analytics-service/internal/consumer/order"
	"github.com/4udiwe/big-bob-pizza/analytics-service/internal/database"
	order_event_repository "github.com/4udiwe/big-bob-pizza/analytics-service/internal/repository/order_event"
	"github.com/4udiwe/big-bob-pizza/analytics-service/internal/service/analytics"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/httpserver"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/kafka"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type App struct {
	cfg       *config.Config
	interrupt chan os.Signal

	// DB
	postgres *postgres.Postgres

	// Echo
	echoHandler *echo.Echo

	// Repositories
	orderEventRepo *order_event_repository.Repository

	// Services
	analyticsService *analytics.Service

	// Consumer
	orderConsumer *consumer_order.Consumer
}

func New(configPath string) *App {
	cfg, err := config.New(configPath)
	if err != nil {
		log.Fatalf("app - New - config.New: %v", err)
	}

	initLogger(cfg.Log.Level)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	return &App{
		cfg:       cfg,
		interrupt: interrupt,
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

	// Prometheus metrics
	metrics := analytics.NewMetrics()

	// Инициализируем сервис с метриками
	app.analyticsService = analytics.NewService(
		app.OrderEventRepo(),
		metrics,
	)

	// Consumer для order.events
	orderKafkaConsumer := kafka.NewConsumer(app.cfg.Kafka.Brokers)

	app.orderConsumer = consumer_order.New(
		app.analyticsService,
		orderKafkaConsumer,
		app.cfg.Kafka.Topics.OrderEvents,
		app.cfg.Kafka.Consumer.GroupID,
	)

	// App server
	log.Info("Starting app server...")
	httpServer := httpserver.New(app.EchoHandler(), httpserver.Port(app.cfg.HTTP.Port))
	httpServer.Start()
	log.Debugf("Server port: %s", app.cfg.HTTP.Port)

	// Run consumer
	app.orderConsumer.Run(ctx)

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
