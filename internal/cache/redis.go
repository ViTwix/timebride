package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"

	"timebride/internal/errors"
)

// Cache визначає інтерфейс для кешування
type Cache interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}

// RedisCache реалізує інтерфейс Cache використовуючи Redis
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache створює новий екземпляр RedisCache
func NewRedisCache(addr, password string, db int) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Перевіряємо підключення
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, errors.NewInternalError("failed to connect to redis", err)
	}

	return &RedisCache{client: client}, nil
}

// Get отримує значення з кешу
func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return errors.NewNotFoundError("cache key not found")
	}
	if err != nil {
		return errors.NewInternalError("failed to get from cache", err)
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return errors.NewInternalError("failed to unmarshal cache value", err)
	}

	return nil
}

// Set зберігає значення в кеш
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return errors.NewInternalError("failed to marshal cache value", err)
	}

	if err := c.client.Set(ctx, key, data, expiration).Err(); err != nil {
		return errors.NewInternalError("failed to set cache", err)
	}

	return nil
}

// Delete видаляє значення з кешу
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	if err := c.client.Del(ctx, key).Err(); err != nil {
		return errors.NewInternalError("failed to delete from cache", err)
	}
	return nil
}

// Exists перевіряє чи існує ключ в кеші
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	val, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, errors.NewInternalError("failed to check cache key", err)
	}
	return val > 0, nil
}

// Close закриває з'єднання з Redis
func (c *RedisCache) Close() error {
	return c.client.Close()
}
