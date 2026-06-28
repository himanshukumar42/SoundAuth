package cache

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

type SessionCache interface {
	SetSession(ctx context.Context, key string, value string, ttl time.Duration) error
	GetSession(ctx context.Context, key string) (string, error)
	DeleteSession(ctx context.Context, key string) error
}

type CacheItem struct {
	Value     string
	ExpiresAt time.Time
}

func NewCacheItem(value string, ttl time.Duration) *CacheItem {
	return &CacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

type RedisCache struct {
	mu   sync.RWMutex
	data map[string]CacheItem
}

func NewRedisCache() *RedisCache {
	return &RedisCache{
		data: make(map[string]CacheItem),
	}
}

func (rc *RedisCache) Get(ctx context.Context, key string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	rc.mu.RLock()
	defer rc.mu.RUnlock()

	cacheItem, ok := rc.data[key]
	if !ok {
		return "", fmt.Errorf("key %s not found", key)
	}

	if time.Now().After(cacheItem.ExpiresAt) {
		rc.mu.Lock()
		delete(rc.data, key)
		rc.mu.Unlock()
		return "", fmt.Errorf("key %s expired", key)
	}

	return cacheItem.Value, nil
}

func (rc *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	rc.mu.Lock()
	defer rc.mu.Unlock()

	cacheItem := NewCacheItem(value, ttl)
	rc.data[key] = *cacheItem
	return nil
}

func (rc *RedisCache) Delete(ctx context.Context, key string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if _, exists := rc.data[key]; !exists {
		return fmt.Errorf("key %s does not exists", key)
	}

	delete(rc.data, key)
	return nil
}

func (rc *RedisCache) StartCleaner(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			now := time.Now()

			rc.mu.Lock()
			for key, item := range rc.data {
				if now.After(item.ExpiresAt) {
					delete(rc.data, key)
				}
			}
			rc.mu.Unlock()
		}
	}()
}
