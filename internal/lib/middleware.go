package lib

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/himanshukumar42/soundauth/internal/models"
)

type AuthenticationMiddleware struct {
	tenantBuckets map[string]*TokenBucket
	userBuckets   map[string]*TokenBucket
	mu            sync.Mutex
}

func NewAuthenticationMiddleware() *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		tenantBuckets: make(map[string]*TokenBucket),
		userBuckets:   make(map[string]*TokenBucket),
	}
}

func (am *AuthenticationMiddleware) GetOrCreateTenantBucket(tenant string) *TokenBucket {
	am.mu.Lock()
	defer am.mu.Unlock()

	bucket, ok := am.tenantBuckets[tenant]
	if ok {
		return bucket
	}

	bucket = NewTokenBucket(100, 20)
	am.tenantBuckets[tenant] = bucket

	return bucket
}

func (am *AuthenticationMiddleware) GetOrCreateUserBucket(userID string) *TokenBucket {
	am.mu.Lock()
	defer am.mu.Unlock()

	bucket, ok := am.userBuckets[userID]
	if ok {
		return bucket
	}

	bucket = NewTokenBucket(10, 5)

	am.userBuckets[userID] = bucket

	return bucket
}

func (am *AuthenticationMiddleware) Process(ctx context.Context, req models.AuthRequest) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	tenant := am.GetOrCreateTenantBucket(req.TenantID)

	if !tenant.Allow() {
		return fmt.Errorf(
			"tenant %s exceeded rate limit", req.TenantID,
		)
	}

	user := am.GetOrCreateUserBucket(req.TenantID + ":" + req.DeviceID)

	if !user.Allow() {
		return fmt.Errorf(
			"user exceeded rate limit",
		)
	}

	log.Printf("[Middleware] Request Accepted Tenant=%s Device=%s", req.TenantID, req.DeviceID)

	return nil
}
