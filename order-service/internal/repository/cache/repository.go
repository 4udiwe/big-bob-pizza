package cache_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/4udiwe/big-bob-pizza/order-service/internal/entity"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/redis"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const (
	orderKey        = "order:%s"
	activeOrdersKey = "orders:active"
	userActiveKey   = "orders:user:%s:active"
	statusOrdersKey = "orders:status:%s"

	orderTTL = 24 * time.Hour
)

type CacheOrderRepository struct {
	client *redis.Redis
}

func NewCacheOrderRepository(client *redis.Redis) *CacheOrderRepository {
	return &CacheOrderRepository{
		client: client,
	}
}

// build keys
func keyOrder(id uuid.UUID) string {
	return fmt.Sprintf(orderKey, id.String())
}

func keyUserActive(userID uuid.UUID) string {
	return fmt.Sprintf(userActiveKey, userID.String())
}

func keyStatus(status string) string {
	return fmt.Sprintf(statusOrdersKey, status)
}

func (r *CacheOrderRepository) Save(ctx context.Context, ord *entity.Order) error {
	b, err := json.Marshal(ord)
	if err != nil {
		return fmt.Errorf("cache order repo - marshal order: %w", err)
	}

	key := keyOrder(ord.ID)

	if err := r.client.Set(ctx, key, b, orderTTL); err != nil {
		return fmt.Errorf("cache order repo - set order: %w", err)
	}

	log.Infof("redis: saved order %s to cache", ord.ID)
	return nil
}

func (r *CacheOrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	key := keyOrder(id)

	data, err := r.client.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("cache order repo - get order: %w", err)
	}

	if data == "" {
		return nil, nil // not found in cache
	}

	var ord entity.Order
	if err := json.Unmarshal([]byte(data), &ord); err != nil {
		return nil, fmt.Errorf("cache order repo - unmarshal order: %w", err)
	}

	log.Debugf("redis: loaded order %s from cache", id)
	return &ord, nil
}

func (r *CacheOrderRepository) AddToActive(ctx context.Context, ord *entity.Order) error {
	if err := r.client.AddToSet(ctx, activeOrdersKey, ord.ID.String()); err != nil {
		return fmt.Errorf("cache order repo - add to active set: %w", err)
	}

	log.Infof("redis: added to active %s", ord.ID)
	return nil
}

func (r *CacheOrderRepository) RemoveFromActive(ctx context.Context, id uuid.UUID) error {
	if err := r.client.RemoveFromSet(ctx, activeOrdersKey, id.String()); err != nil {
		return fmt.Errorf("cache order repo - remove from active set: %w", err)
	}

	log.Infof("redis: removed from active %s", id)
	return nil
}

func (r *CacheOrderRepository) AddUserActive(ctx context.Context, userID uuid.UUID, ordID uuid.UUID) error {
	if err := r.client.AddToSet(ctx, keyUserActive(userID), ordID.String()); err != nil {
		return fmt.Errorf("cache order repo - add user active: %w", err)
	}

	return nil
}

func (r *CacheOrderRepository) RemoveUserActive(ctx context.Context, userID uuid.UUID, ordID uuid.UUID) error {
	if err := r.client.RemoveFromSet(ctx, keyUserActive(userID), ordID.String()); err != nil {
		return fmt.Errorf("cache order repo - remove user active: %w", err)
	}

	return nil
}

func (r *CacheOrderRepository) AddToStatus(ctx context.Context, status string, ordID uuid.UUID) error {
	if err := r.client.AddToSet(ctx, keyStatus(status), ordID.String()); err != nil {
		return fmt.Errorf("cache order repo - add to status set: %w", err)
	}
	return nil
}

func (r *CacheOrderRepository) RemoveFromStatus(ctx context.Context, status string, ordID uuid.UUID) error {
	if err := r.client.RemoveFromSet(ctx, keyStatus(status), ordID.String()); err != nil {
		return fmt.Errorf("cache order repo - remove from status set: %w", err)
	}
	return nil
}

func (r *CacheOrderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.client.Delete(ctx, keyOrder(id)); err != nil {
		return fmt.Errorf("cache order repo - delete order: %w", err)
	}

	log.Infof("redis: removed order %s from cache", id)
	return nil
}

func (r *CacheOrderRepository) GetActiveOrders(ctx context.Context) ([]string, error) {
	ids, err := r.client.GetSetMembers(ctx, activeOrdersKey)
	if err != nil {
		return nil, fmt.Errorf("cache order repo - get active members: %w", err)
	}

	return ids, nil
}

func (r *CacheOrderRepository) GetUserActiveOrders(ctx context.Context, userID uuid.UUID) ([]string, error) {
	ids, err := r.client.GetSetMembers(ctx, keyUserActive(userID))
	if err != nil {
		return nil, fmt.Errorf("cache order repo - get user active orders: %w", err)
	}

	return ids, nil
}

func (r *CacheOrderRepository) GetByStatus(ctx context.Context, status string) ([]string, error) {
	ids, err := r.client.GetSetMembers(ctx, keyStatus(status))
	if err != nil {
		return nil, fmt.Errorf("cache order repo - get by status: %w", err)
	}

	return ids, nil
}
