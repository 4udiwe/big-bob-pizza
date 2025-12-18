package order_cache_repository

import (
	"context"
	"errors"
	"time"

	"github.com/4udiwe/big-bob-pizza/payment-service/internal/entity"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *Repository {
	return &Repository{Postgres: pg}
}

func (r *Repository) Save(ctx context.Context, orderInfo entity.OrderInfo) error {
	logrus.Infof("OrderCacheRepository.Save: orderID=%s", orderInfo.OrderID)

	// Заказы доступны для оплаты в течение 30 минут
	expiresAt := time.Now().Add(30 * time.Minute)

	query, args, _ := r.Builder.
		Insert("order_cache").
		Columns("order_id", "user_id", "total_price", "created_at", "expires_at").
		Values(orderInfo.OrderID, orderInfo.UserID, orderInfo.TotalPrice, orderInfo.CreatedAt, expiresAt).
		Suffix("ON CONFLICT (order_id) DO UPDATE SET expires_at = EXCLUDED.expires_at").
		ToSql()

	_, err := r.GetTxManager(ctx).Exec(ctx, query, args...)
	if err != nil {
		logrus.Errorf("OrderCacheRepository.Save: error: %v", err)
		return err
	}

	return nil
}

func (r *Repository) GetByOrderID(ctx context.Context, orderID uuid.UUID) (entity.OrderInfo, error) {
	query := `
		SELECT order_id, user_id, total_price, created_at
		FROM order_cache
		WHERE order_id = $1 AND expires_at > NOW()
	`

	row := r.GetTxManager(ctx).QueryRow(ctx, query, orderID)
	var orderInfo entity.OrderInfo
	if err := row.Scan(&orderInfo.OrderID, &orderInfo.UserID, &orderInfo.TotalPrice, &orderInfo.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.OrderInfo{}, ErrOrderNotFound
		}
		return entity.OrderInfo{}, err
	}

	return orderInfo, nil
}

func (r *Repository) Delete(ctx context.Context, orderID uuid.UUID) error {
	query, args, _ := r.Builder.
		Delete("order_cache").
		Where("order_id = ?", orderID).
		ToSql()

	_, err := r.GetTxManager(ctx).Exec(ctx, query, args...)
	if err != nil {
		logrus.Errorf("OrderCacheRepository.Delete: error: %v", err)
		return err
	}

	return nil
}
