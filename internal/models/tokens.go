package models

import (
	"context"
	"time"
)

const (
	DefaultTokenTTL = time.Hour
)

type JWTClaims struct {
	UserID    string
	Email     string
	TenantID  string
	SessionID string
	Roles     []string
	Scopes    []string

	Issuer    string
	Audience  string
	ExpiresAt time.Time
	IssuedAt  time.Time
	NotBefore time.Time
}

type TokenManager interface {
	GenerateToken(ctx context.Context, req GenerateTokenRequest) (*TokenPair, error)
	VerifyToken(ctx context.Context, token string) (*JWTClaims, error)
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int64
}

type GenerateTokenRequest struct {
	UserID    string
	Email     string
	TenantID  string
	SessionID string
	Roles     []Role
	Scopes    []Scope
}
