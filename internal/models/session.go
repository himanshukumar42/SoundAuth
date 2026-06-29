package models

import (
	"context"
	"time"
)

type SessionStatus string

const (
	SessionActive  SessionStatus = "active"
	SessionRevoked SessionStatus = "revoked"
	SessionExpired SessionStatus = "expired"
)

type SessionManager interface {
	CreateSession(ctx context.Context, req CreateSessionRequest) (*Session, error)
	GetSession(ctx context.Context, sessionID string) (*Session, error)
	UpdateSession(ctx context.Context, req UpdateSessionRequest) error
	RevokeSession(ctx context.Context, sessionID string) error
	DeleteSession(ctx context.Context, sessionID string) error
}

type SessionStore interface {
	Save(ctx context.Context, session *Session) error
	Get(ctx context.Context, id string) (*Session, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, session *Session) error
}

type Session struct {
	ID       string
	UserID   string
	TenantID string

	DeviceID   string
	DeviceName string
	UserAgent  string
	IPAddress  string
	Status     SessionStatus
	CreatedAt  time.Time
	LastSeenAt time.Time
	ExpiresAt  time.Time
}

type CreateSessionRequest struct {
	UserID     string
	TenantID   string
	DeviceID   string
	DeviceName string
	UserAgent  string
	IPAddress  string
}

type UpdateSessionRequest struct {
	SessionID  string
	LastSeenAt time.Time
	UserAgent  string
	IPAddress  string
}
