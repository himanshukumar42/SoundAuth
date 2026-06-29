package provider

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/himanshukumar42/soundauth/internal/models"
)

type AuthenticationProvider interface {
	Name() models.Provider
	Authenticate(ctx context.Context, request models.AuthRequest) (*models.ProviderResponse, error)
}

// Factory Pattern
type ProviderFactory struct {
	mu        sync.RWMutex
	providers map[models.Provider]AuthenticationProvider
}

func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{
		providers: make(map[models.Provider]AuthenticationProvider),
	}
}

func (pf *ProviderFactory) Register(provider AuthenticationProvider) {
	pf.mu.Lock()
	defer pf.mu.Unlock()

	pf.providers[provider.Name()] = provider
}

func (pf *ProviderFactory) Get(name models.Provider) (AuthenticationProvider, error) {
	pf.mu.RLock()
	defer pf.mu.RUnlock()
	providerName := models.Provider(name)
	provider, exists := pf.providers[providerName]
	if !exists {
		return nil, fmt.Errorf("provider %s not registered", name)
	}
	return provider, nil
}

type ResponseAdapter interface {
	Normalize(any) (*models.ProviderResponse, error)
}

// Adapter Pattern
type DefaultAdapter struct{}

func NewDefaultAdapter() *DefaultAdapter {
	return &DefaultAdapter{}
}

func (a *DefaultAdapter) Normalize(response any) (*models.ProviderResponse, error) {
	data, ok := response.(map[string]any)
	if !ok {
		return nil, errors.New("invalid provider response")
	}

	userID, _ := data["user_id"].(string)
	status, _ := data["status"].(bool)

	return &models.ProviderResponse{
		UserID:        userID,
		Authenticated: status,
		Metadata: map[string]string{
			"normalized": "true",
		},
	}, nil
}

// Passkeys Authentication
type PasskeyProvider struct {
	Adapter ResponseAdapter
}

func NewPasskeyProvider(adapter ResponseAdapter) *PasskeyProvider {
	return &PasskeyProvider{
		Adapter: adapter,
	}
}

func (pp *PasskeyProvider) Name() models.Provider {
	return models.ProviderPasskey
}

func (pp *PasskeyProvider) Authenticate(ctx context.Context, req models.AuthRequest) (*models.ProviderResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:

	}
	log.Printf("[Passkey] Authenticating Device %s\n", req.DeviceID)

	// ------------------------------------------
	// Simulate WebAuthn verification
	// In production this would:
	// • Verify challenge
	// • Verify authenticator
	// • Verify signature
	// • Verify RP ID
	// • Verify origin
	// ------------------------------------------

	raw := map[string]any{
		"user_id": "USER-101",
		"status":  true,
		"device":  req.DeviceID,
	}
	return pp.Adapter.Normalize(raw)
}

// Google OAuth Authentication
type GoogleOAuthProvider struct {
	Adapter ResponseAdapter
}

func NewGoogleOAuthProvider(adapter ResponseAdapter) *GoogleOAuthProvider {
	return &GoogleOAuthProvider{
		Adapter: adapter,
	}
}

func (gp *GoogleOAuthProvider) Name() models.Provider {
	return models.ProviderGoogle
}

func (gp *GoogleOAuthProvider) Authenticate(ctx context.Context, req models.AuthRequest) (*models.ProviderResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	log.Printf("[Google OAuth] Validating OAuth token")

	raw := map[string]any{
		"user_id": "USER-101",
		"status":  true,
		"email":   "user@gmail.com",
	}
	return gp.Adapter.Normalize(raw)
}
