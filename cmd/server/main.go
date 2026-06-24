package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type AuthRequest struct {
	TenantID   string
	Provider   string
	Credential string
	DeviceID   string
}

type AuthResponse struct {
	Authenticated bool
	UserID        string
	Provider      string
	Token         string
	ExpiresIn     int64
}

type TenantConfig struct {
	Name       string
	Provider   string
	PublicKey  string
	SecretPath string
}

type User struct {
	ID       string
	Email    string
	Tenant   string
	DeviceID string
}

type AuthenticationProvider interface {
	Name() string
	Authenticate(ctx context.Context, request AuthRequest) (*AuthResponse, error)
}

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*User, error)
}

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

type AuditEvent struct {
	ID        string                 `json:"id"`
	TenantID  string                 `json:"tenant_id"`
	UserID    string                 `json:"user_id"`
	EventType string                 `json:"event_type"`
	Timestamp time.Time              `json:"timestamp"`
	Actor     string                 `json:"actor"` // Email or System
	Success   bool                   `json:"success"`
	IPAddress string                 `json:"ip_address,omitempty"`
	UserAgent string                 `json:"user_agent,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

type QueryFilter struct {
	EventType string
	Success   *bool
	Limit     int
	Offset    int
}

type AuditLogRepository interface {
	LogEvent(ctx context.Context, event *AuditEvent) error
	QueryEvents(ctx context.Context, tenantID string, filter QueryFilter) ([]AuditEvent, error)
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

/