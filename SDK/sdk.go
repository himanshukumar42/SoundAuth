package sdk

import (
	"time"

	"github.com/himanshukumar42/soundauth/internal/auth"
	"github.com/himanshukumar42/soundauth/internal/cache"
	"github.com/himanshukumar42/soundauth/internal/lib"
	"github.com/himanshukumar42/soundauth/internal/provider"
	"github.com/himanshukumar42/soundauth/internal/repository"
	"github.com/himanshukumar42/soundauth/internal/services"
	"github.com/himanshukumar42/soundauth/internal/vault"
	"github.com/himanshukumar42/soundauth/internal/worker"
)

const (
	Secret = "AUTH_APPLICATION_SECRET"
	Issuer = "SOUNDAUTH"
	TTL    = 5 * time.Minute
)

func NewAuthenticationSDK() *auth.AuthService {
	// Infrastructure

	cache := cache.NewRedisCache()
	vault := vault.NewVaultClient()
	userRepo := repository.NewInMemoryUserRepository()

	sessionStore := auth.NewInMemorySessionStore()
	sessionManager := auth.NewDefaultSessionManager(sessionStore)

	tokenManager := auth.NewJWTTokenManager(Secret, Issuer, TTL)
	tenantService := services.NewTenantService(cache, vault)

	pool := worker.NewVerificationPool(
		4, // workers
		2, // concurrent crypto verification
	)

	verifier := services.NewSignatureVerifier(tenantService, pool)

	middleware := lib.NewAuthenticationMiddleware()

	// Factory

	factory := provider.NewProviderFactory()
	adapter := provider.NewDefaultAdapter()

	// PasskeyProvider

	var prd provider.AuthenticationProvider
	prd = provider.NewPasskeyProvider(adapter)
	// prd = lib.NewAuditDecorator(prd)
	// prd = lib.NewMetricsDecorator(prd)
	// prd = lib.NewLoggingDecorator(prd)
	factory.Register(prd)

	return &auth.AuthService{
		Factory:        factory,
		UserRepo:       userRepo,
		Cache:          cache,
		Vault:          vault,
		SessionManager: sessionManager,
		TokenManager:   tokenManager,
		Verifier:       verifier,
		Middleware:     middleware,
	}
}
