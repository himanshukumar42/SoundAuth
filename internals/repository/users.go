package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/himanshukumar42/soundauth/internals/models"
)

type InMemoryUserRepository struct {
	mu    sync.RWMutex
	Users map[string]models.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		Users: map[string]models.User{
			"USER-101": {
				ID:       "USER-101",
				Email:    "passkey@example.com",
				Tenant:   "Google",
				DeviceID: "DEVICE-101",
			},
			"USER-202": {
				ID:       "USER-202",
				Email:    "google@example.com",
				Tenant:   "Google",
				DeviceID: "DEVICE-202",
			},
			"USER-303": {
				ID:       "USER-303",
				Email:    "github@example.com",
				Tenant:   "GitHub",
				DeviceID: "DEVICE-303",
			},
		},
	}
}

func (ur *InMemoryUserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	ur.mu.RLock()
	defer ur.mu.RUnlock()

	user, ok := ur.Users[id]
	if !ok {
		return nil, fmt.Errorf("usre %s not found", id)
	}
	return &user, nil
}
