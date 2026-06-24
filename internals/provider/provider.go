package provider

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/himanshukumar42/soundauth/internals/models"
)

type Provider string

const (
	ProviderPasskey Provider = "Passkey"
	ProviderGoogle  Provider = "GoogleOAuth"
	ProviderGithub  Provider = "GithubOAuth"
	ProviderMagic   Provider = "MagicLink"
	ProviderSAML    Provider = "SAML"
)

type AuthenticationProvider interface {
	Name() string
	Authenticate(ctx context.Context, request models.AuthRequest) (*models.AuthResponse, error)
}

// Factory Pattern
type ProviderFactory struct {
	mu        sync.RWMutex
	providers map[string]AuthenticationProvider
}

func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{
		providers: make(map[string]AuthenticationProvider),
	}
}

func (pf *ProviderFactory) Register(provider AuthenticationProvider) {
	pf.mu.Lock()
	defer pf.mu.Unlock()

	pf.providers[provider.Name()] = provider
}

func (pf *ProviderFactory) Get(name string) (AuthenticationProvider, error) {
	pf.mu.RLock()
	defer pf.mu.RUnlock()
	provider, exists := pf.providers[name]
	if !exists {
		return nil, fmt.Errorf("provider %s not registered", name)
	}
	return provider, nil
}

type ProviderResponse struct {
	UserID        string
	Authenticated bool
	Metadata      map[string]string
}

type ResponseAdapter interface {
	Normalize(any) (*ProviderResponse, error)
}

// Adapter Pattern
type DefaultAdapter struct{}

func NewDefaultAdapter() *DefaultAdapter {
	return &DefaultAdapter{}
}

func (a *DefaultAdapter) Normalize(response any) (*ProviderResponse, error) {
	data, ok := response.(map[string]any)
	if !ok {
		return nil, errors.New("invalid provider response")
	}

	userID, _ := data["user_id"].(string)
	status, _ := data["status"].(bool)

	return &ProviderResponse{
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

func (pp *PasskeyProvider) Name() Provider {
	return ProviderPasskey
}

func (pp *PasskeyProvider) Authenticate(ctx context.Context, req models.AuthRequest) (*ProviderResponse, error) {
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

func (gp *GoogleOAuthProvider) Name() Provider {
	return ProviderGoogle
}

func (gp *GoogleOAuthProvider) Authenticate(ctx context.Context, req models.AuthRequest) (*ProviderResponse, error) {
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
