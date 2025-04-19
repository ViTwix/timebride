package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Cache interface defines the methods that all cache implementations must provide
type Cache interface {
	// Get retrieves a value from the cache
	Get(ctx context.Context, key string) (interface{}, error)

	// Set stores a value in the cache with an optional expiration time
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error

	// Delete removes a value from the cache
	Delete(ctx context.Context, key string) error

	// Exists checks if a key exists in the cache
	Exists(ctx context.Context, key string) (bool, error)

	// Clear removes all items from the cache
	Clear(ctx context.Context) error

	// Close closes the cache connection
	Close() error
}

// RedisCache implements the Cache interface using Redis
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis cache client
func NewRedisCache(host string, port int, password string, db int, poolSize int, idleTimeout time.Duration) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		Password:     password,
		DB:           db,
		PoolSize:     poolSize,
		IdleTimeout:  idleTimeout,
		MinIdleConns: 2,
	})

	// Test the connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCache{client: client}, nil
}

// Get retrieves a value from Redis
func (c *RedisCache) Get(ctx context.Context, key string) (interface{}, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, err
	}
	return result, nil
}

// Set stores a value in Redis
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, expiration).Err()
}

// Delete removes a value from Redis
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Exists checks if a key exists in Redis
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// Clear removes all items from Redis
func (c *RedisCache) Clear(ctx context.Context) error {
	return c.client.FlushDB(ctx).Err()
}

// Close closes the Redis connection
func (c *RedisCache) Close() error {
	return c.client.Close()
}

// InMemoryCache implements the Cache interface using an in-memory map
type InMemoryCache struct {
	store   map[string]*cacheItem
	maxSize int
	ttl     time.Duration
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

// NewInMemoryCache creates a new in-memory cache
func NewInMemoryCache(maxSize int, ttl time.Duration) *InMemoryCache {
	cache := &InMemoryCache{
		store:   make(map[string]*cacheItem),
		maxSize: maxSize,
		ttl:     ttl,
	}

	// Start cleanup goroutine
	go cache.cleanup()

	return cache
}

// cleanup removes expired items from the cache
func (c *InMemoryCache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		for key, item := range c.store {
			if now.After(item.expiration) {
				delete(c.store, key)
			}
		}
	}
}

// Get retrieves a value from the in-memory cache
func (c *InMemoryCache) Get(ctx context.Context, key string) (interface{}, error) {
	item, exists := c.store[key]
	if !exists {
		return nil, nil
	}

	if time.Now().After(item.expiration) {
		delete(c.store, key)
		return nil, nil
	}

	return item.value, nil
}

// Set stores a value in the in-memory cache
func (c *InMemoryCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	// If max size is reached, remove the oldest item
	if len(c.store) >= c.maxSize {
		var oldestKey string
		var oldestTime time.Time

		for k, v := range c.store {
			if oldestKey == "" || v.expiration.Before(oldestTime) {
				oldestKey = k
				oldestTime = v.expiration
			}
		}

		if oldestKey != "" {
			delete(c.store, oldestKey)
		}
	}

	// Set default expiration if not provided
	if expiration == 0 {
		expiration = c.ttl
	}

	c.store[key] = &cacheItem{
		value:      value,
		expiration: time.Now().Add(expiration),
	}

	return nil
}

// Delete removes a value from the in-memory cache
func (c *InMemoryCache) Delete(ctx context.Context, key string) error {
	delete(c.store, key)
	return nil
}

// Exists checks if a key exists in the in-memory cache
func (c *InMemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	item, exists := c.store[key]
	if !exists {
		return false, nil
	}

	if time.Now().After(item.expiration) {
		delete(c.store, key)
		return false, nil
	}

	return true, nil
}

// Clear removes all items from the in-memory cache
func (c *InMemoryCache) Clear(ctx context.Context) error {
	c.store = make(map[string]*cacheItem)
	return nil
}

// Close closes the in-memory cache
func (c *InMemoryCache) Close() error {
	// Nothing to close for in-memory cache
	return nil
}

// CacheFactory creates a cache based on the provided configuration
func CacheFactory(cacheType string, config map[string]interface{}) (Cache, error) {
	switch cacheType {
	case "redis":
		host := config["host"].(string)
		port := config["port"].(int)
		password := config["password"].(string)
		db := config["db"].(int)
		poolSize := config["pool_size"].(int)
		idleTimeout := time.Duration(config["idle_timeout"].(int)) * time.Second

		return NewRedisCache(host, port, password, db, poolSize, idleTimeout)

	case "in-memory":
		maxSize := config["max_size"].(int)
		ttl := time.Duration(config["ttl"].(int)) * time.Second

		return NewInMemoryCache(maxSize, ttl), nil

	default:
		return nil, fmt.Errorf("unsupported cache type: %s", cacheType)
	}
}
