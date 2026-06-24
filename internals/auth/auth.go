package auth

import (
	"context"
	"log"
)

// Dependency Injection

type AuthService struct {
	Factory      *ProviderFactory
	UserRepo     UserRepository
	Cache        Cache
	SessionCache SessionCache
	TokenManager TokenManager
	AuditLog     AuditLogRepository
}

func NewAuthService(factory *ProviderFactory, userRepo UserRepository, cache Cache, sessionCache SessionCache, tokenManager TokenManager, auditLog AuditLogRepository) *AuthService {
	return &AuthService{
		Factory:      factory,
		UserRepo:     userRepo,
		Cache:        cache,
		SessionCache: sessionCache,
		TokenManager: tokenManager,
		AuditLog:     auditLog,
	}
}

func (as *AuthService) Authenticate(ctx context.Context, req AuthRequest) (*AuthResponse, error) {
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

	return &AuthResponse{
		Authenticated: true,
		UserID:        user.ID,
		Provider:      req.Provider,
		Token:         token,
		ExpiresIn:     int64(DefaultTokenTTL.Seconds()),
	}, nil
}
