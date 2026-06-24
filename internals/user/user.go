package user

import "context"

type User struct {
	ID       string
	Email    string
	Tenant   string
	DeviceID string
}

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*User, error)
}
