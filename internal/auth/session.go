package auth

import (
	"context"
	"time"

	"github.com/himanshukumar42/soundauth/internal/models"
)

type DefaultSessionManager struct {
	store models.SessionStore
}

func NewDefaultSessionManager(store models.SessionStore) *DefaultSessionManager {
	return &DefaultSessionManager{
		store: store,
	}
}

func (dsm *DefaultSessionManager) CreateSession(ctx context.Context, req models.CreateSessionRequest) (*models.Session, error) {
	return &models.Session{
		ID:         "sess_01HZX8A9K3",
		UserID:     req.UserID,
		TenantID:   req.TenantID,
		DeviceID:   req.DeviceID,
		DeviceName: req.DeviceName,
		UserAgent:  req.UserAgent,
		IPAddress:  req.IPAddress,
		Status:     models.SessionActive,
		CreatedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(time.Hour),
	}, nil
}

func (dsm *DefaultSessionManager) GetSession(ctx context.Context, sessionID string) (*models.Session, error) {
	session, err := dsm.store.Get(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (dsm *DefaultSessionManager) UpdateSession(ctx context.Context, req models.UpdateSessionRequest) error {
	session, err := dsm.store.Get(ctx, req.SessionID)
	if err != nil {
		return err
	}
	err = dsm.store.Update(ctx, session)
	return err
}

func (dsm *DefaultSessionManager) RevokeSession(ctx context.Context, sessionID string) error {
	_, err := dsm.store.Get(ctx, sessionID)
	if err != nil {
		return err
	}
	err = dsm.store.Delete(ctx, sessionID)
	return err
}

func (dsm *DefaultSessionManager) DeleteSession(ctx context.Context, sessionID string) error {
	err := dsm.store.Delete(ctx, sessionID)
	return err
}
