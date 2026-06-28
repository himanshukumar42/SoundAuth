package models

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
