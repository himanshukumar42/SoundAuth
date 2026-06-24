package auth

import (
	"context"
	"log"
	"time"

	"github.com/himanshukumar42/soundauth/internals/audit"
	"github.com/himanshukumar42/soundauth/internals/cache"
	"github.com/himanshukumar42/soundauth/internals/models"
	"github.com/himanshukumar42/soundauth/internals/provider"
	"github.com/himanshukumar42/soundauth/internals/user"
)

const (
	DefaultTokenTTL = time.Hour
)

type Claims struct {
	UserID    string
	Audience  string
	ExpiresAt int64
	IssuedAt  int64
}

type TokenManager interface {
	GenerateToken(ctx context.Context, userId string) (string, error)
	VerifyToken(ctx context.Context, token string) (*Claims, error)
}

// Dependency Injection
type AuthService struct {
	Factory      *provider.ProviderFactory
	UserRepo     user.UserRepository
	Cache        cache.Cache
	SessionCache cache.SessionCache
	TokenManager TokenManager
	AuditLog     audit.AuditLogRepository
}

func NewAuthService(factory *provider.ProviderFactory, userRepo user.UserRepository, cache cache.Cache, sessionCache cache.SessionCache, tokenManager TokenManager, auditLog audit.AuditLogRepository) *AuthService {
	return &AuthService{
		Factory:      factory,
		UserRepo:     userRepo,
		Cache:        cache,
		SessionCache: sessionCache,
		TokenManager: tokenManager,
		AuditLog:     auditLog,
	}
}

func (as *AuthService) Authenticate(ctx context.Context, req models.AuthRequest) (*models.AuthResponse, error) {
	provider, err := as.Factory.Get(req.Provider)
	if err != nil {
		return nil, err
	}

	log.Printf("[Auth] Tenant=%s Provider=%s", req.TenantID, req.Provider)

	authResponse, err := provider.Authenticate(ctx, req)
	if err != nil {
		return nil, err
	}

	if !authResponse.Authenticated {
		return nil, err
	}

	user, err := as.UserRepo.FindByID(ctx, authResponse.UserID)
	if err != nil {
		return nil, err
	}

	token, err := as.TokenManager.GenerateToken(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Authenticated: true,
		UserID:        user.ID,
		Provider:      req.Provider,
		Token:         token,
		ExpiresIn:     int64(DefaultTokenTTL.Seconds()),
	}, nil
}
