package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

const (
	defaultConnTimeout  = time.Second
	defaultConnAttempts = 10
)

type Redis struct {
	Client       *redis.Client
	connAttempts int
	connTimeout  time.Duration
	keyPrefix    string
}

func New(addr, password string, db int, opts ...Option) (*Redis, error) {
	r := &Redis{
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
		keyPrefix:    "",
	}

	// apply custom options
	for _, opt := range opts {
		opt(r)
	}

	r.Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	var err error
	for r.connAttempts > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), r.connTimeout)
		defer cancel()

		_, err = r.Client.Ping(ctx).Result()
		if err == nil {
			break
		}

		log.Infof("Redis is trying to connect, attempts left: %d", r.connAttempts)
		r.connAttempts--
		time.Sleep(r.connTimeout)
	}

	if err != nil {
		return nil, fmt.Errorf("redis - NewRedis - connAttempts == 0: %w", err)
	}

	return r, nil
}

func (r *Redis) Close() error {
	if r.Client != nil {
		return r.Client.Close()
	}
	return nil
}

// Key helper — applies prefix if set
func (r *Redis) key(k string) string {
	if r.keyPrefix == "" {
		return k
	}
	return r.keyPrefix + ":" + k
}

// Common Redis helpers

func (r *Redis) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.Client.Set(ctx, r.key(key), value, ttl).Err()
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, r.key(key)).Result()
}

func (r *Redis) Exists(ctx context.Context, key string) (bool, error) {
	n, err := r.Client.Exists(ctx, r.key(key)).Result()
	return n == 1, err
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, r.key(key)).Err()
}

func (r *Redis) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return r.Client.Expire(ctx, r.key(key), ttl).Err()
}

// Set operations

func (r *Redis) AddToSet(ctx context.Context, key string, members ...interface{}) error {
	return r.Client.SAdd(ctx, r.key(key), members...).Err()
}

func (r *Redis) RemoveFromSet(ctx context.Context, key string, members ...interface{}) error {
	return r.Client.SRem(ctx, r.key(key), members...).Err()
}

func (r *Redis) GetSetMembers(ctx context.Context, key string) ([]string, error) {
	return r.Client.SMembers(ctx, r.key(key)).Result()
}
