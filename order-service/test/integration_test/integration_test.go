//go:build integration

package integration_test

import (
	"context"
	"os"
	"testing"

	"github.com/4udiwe/big-bob-pizza/order-service/internal/database"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/postgres"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/redis"
	log "github.com/sirupsen/logrus"
)

const (
	defaultAttempts = 20
)

var (
	testPostgres *postgres.Postgres
	testRedis    *redis.Redis
)

func TestMain(m *testing.M) {
	postgresURL := os.Getenv("POSTGRES_URL")
	if postgresURL == "" {
		postgresURL = "postgres://test_user:test_password@localhost:5432/test_db"
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	var err error
	testPostgres, err = postgres.New(postgresURL, postgres.ConnAttempts(5))
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer testPostgres.Close()

	// Run migrations
	if err := database.RunMigrations(context.Background(), testPostgres.Pool); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	testRedis, err = redis.New(redisAddr, "", 0)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	defer testRedis.Close()

	log.Info("integration tests: database connections established")
	os.Exit(m.Run())
}


