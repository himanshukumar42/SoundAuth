package auth

import (
	"context"
	"log"

	"github.com/himanshukumar42/soundauth/internal/audit"
	"github.com/himanshukumar42/soundauth/internal/cache"
	"github.com/himanshukumar42/soundauth/internal/lib"
	"github.com/himanshukumar42/soundauth/internal/models"
	"github.com/himanshukumar42/soundauth/internal/provider"
	"github.com/himanshukumar42/soundauth/internal/repository"
	"github.com/himanshukumar42/soundauth/internal/services"
	"github.com/himanshukumar42/soundauth/internal/vault"
)

// Dependency Injection
type AuthService struct {
	Middleware     *lib.AuthenticationMiddleware
	Factory        *provider.ProviderFactory
	Verifier       *services.SignatureVerifier
	UserRepo       *repository.InMemoryUserRepository
	Cache          cache.Cache
	SessionManager models.SessionManager
	Vault          *vault.VaultClient
	TokenManager   models.TokenManager
	AuditLog       *audit.AuditLogRepository
}

func NewAuthService(middleware *lib.AuthenticationMiddleware, factory *provider.ProviderFactory, verifier *services.SignatureVerifier, userRepo *repository.InMemoryUserRepository, cache cache.Cache, sessionManager models.SessionManager, vault *vault.VaultClient, tokenManager models.TokenManager, auditLog *audit.AuditLogRepository) *AuthService {
	return &AuthService{
		Middleware:     middleware,
		Factory:        factory,
		Verifier:       verifier,
		UserRepo:       userRepo,
		Cache:          cache,
		SessionManager: sessionManager,
		Vault:          vault,
		TokenManager:   tokenManager,
		AuditLog:       auditLog,
	}
}

func (as *AuthService) Authenticate(ctx context.Context, req models.AuthRequest) (*models.AuthResponse, error) {

	// 1. Middleware
	if err := as.Middleware.Process(ctx, req); err != nil {
		return nil, err
	}

	// 2. Factory => Authentication Provider

	provider, err := as.Factory.Get(req.Provider)
	if err != nil {
		return nil, err
	}

	log.Printf("[Auth] Tenant=%s Provider=%s", req.TenantID, req.Provider)

	// 3. Authentication
	authResponse, err := provider.Authenticate(ctx, req)
	if err != nil {
		return nil, err
	}

	if !authResponse.Authenticated {
		return nil, err
	}

	// 4. Signature Verification
	if err := as.Verifier.Verify(ctx, req); err != nil {
		return nil, err
	}

	// 5. User Lookup
	user, err := as.UserRepo.FindByID(ctx, authResponse.UserID)
	if err != nil {
		return nil, err
	}

	// Session Start
	sessionReq := models.CreateSessionRequest{
		UserID:    user.ID,
		TenantID:  req.TenantID,
		DeviceID:  req.DeviceID,
		UserAgent: req.UserAgent,
		IPAddress: req.IPAddress,
	}
	session, err := as.SessionManager.CreateSession(ctx, sessionReq)
	if err != nil {
		return nil, err
	}

	// 6. JWT Generation
	tokenReq := models.GenerateTokenRequest{
		UserID:    user.ID,
		Email:     user.Email,
		TenantID:  req.TenantID,
		SessionID: session.ID, // will create in a while
		Roles:     user.Roles,
		Scopes:    user.Scopes,
	}
	token, err := as.TokenManager.GenerateToken(ctx, tokenReq)
	if err != nil {
		return nil, err
	}

	// 7. Build Response
	return &models.AuthResponse{
		Authenticated: true,
		UserID:        user.ID,
		Provider:      req.Provider,
		AccessToken:   token.AccessToken,
		RefreshToken:  token.RefreshToken,
		ExpiresIn:     token.ExpiresIn,
	}, nil
}
