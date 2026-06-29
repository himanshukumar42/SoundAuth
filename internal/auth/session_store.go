package auth

import (
	"context"
	"fmt"
	"sync"

	"github.com/himanshukumar42/soundauth/internal/models"
)

type InMemorySessionStore struct {
	sessions map[string]*models.Session
	mu       sync.RWMutex
}

func NewInMemorySessionStore() *InMemorySessionStore {
	return &InMemorySessionStore{
		sessions: make(map[string]*models.Session),
	}
}

func (imss *InMemorySessionStore) Save(ctx context.Context, session *models.Session) error {
	imss.mu.Lock()
	defer imss.mu.Unlock()

	imss.sessions[session.ID] = session
	return nil
}

func (imss *InMemorySessionStore) Get(ctx context.Context, sessionId string) (*models.Session, error) {
	imss.mu.RLock()
	defer imss.mu.RUnlock()

	session, ok := imss.sessions[sessionId]
	if !ok {
		return nil, fmt.Errorf("session %v does not exists \n", sessionId)
	}
	return session, nil
}

func (imss *InMemorySessionStore) Update(ctx context.Context, session *models.Session) error {
	imss.mu.Lock()
	defer imss.mu.Unlock()

	if _, ok := imss.sessions[session.ID]; !ok {
		return fmt.Errorf("session %v does not exists \n", session.ID)
	}
	imss.sessions[session.ID] = session
	return nil
}

func (imss *InMemorySessionStore) Delete(ctx context.Context, sessionId string) error {
	imss.mu.Lock()
	defer imss.mu.Unlock()

	if _, exists := imss.sessions[sessionId]; !exists {
		return fmt.Errorf("session %v does not exists \n", sessionId)
	}
	delete(imss.sessions, sessionId)
	return nil
}
