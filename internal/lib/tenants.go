package lib

import (
	"context"
	"log"
	"time"

	"github.com/himanshukumar42/soundauth/internal/cache"
	"github.com/himanshukumar42/soundauth/internal/vault"
)

const (
	RetryAttempts   = 3
	RetryTime       = 200 * time.Millisecond
	DefaultTokenTTL = 10 * time.Minute
)

type TenantService struct {
	cache cache.Cache
	vault *vault.VaultClient
}

func NewTenantService(cache cache.Cache, vault *vault.VaultClient) *TenantService {
	return &TenantService{
		cache: cache,
		vault: vault,
	}
}

func (t *TenantService) GetSigninKey(ctx context.Context, tenant string) (string, error) {
	cacheKey := tenant + ":signing-key"

	value, err := t.cache.Get(ctx, cacheKey)
	if err == nil {
		log.Printf("[Cache] HIT tenant - %s\n", tenant)

		return value, nil
	}

	log.Printf("[Cache] Miss tenant - %s\n", tenant)

	var secret string

	err = Retry(ctx, RetryAttempts, RetryTime, func() error {
		var e error

		secret, e = t.vault.GetSecret(ctx, tenant+"/signing-key")

		return e
	})

	if err != nil {
		return "", err
	}

	_ = t.cache.Set(ctx, cacheKey, secret, DefaultTokenTTL)

	return secret, nil
}
