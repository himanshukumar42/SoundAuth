package cache

import (
	"context"
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
