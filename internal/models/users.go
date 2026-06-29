package models

type Role string

type Scope string

const (
	Admin   Role = "Admin"
	Manager Role = "Manager"
	General Role = "General"
)

const (
	UserRead   Scope = "users:read"
	UserWrite  Scope = "users:write"
	UserDelete Scope = "users:delete"
)

type User struct {
	ID       string
	Email    string
	Tenant   string
	Roles    []Role
	Scopes   []Scope
	DeviceID string
}
