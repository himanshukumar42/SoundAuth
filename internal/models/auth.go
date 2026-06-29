package models

type AuthRequest struct {
	TenantID   string
	Provider   Provider
	Credential string
	DeviceID   string
	UserAgent  string
	IPAddress  string
}

type AuthResponse struct {
	Authenticated bool
	UserID        string
	Provider      Provider
	AccessToken   string
	RefreshToken  string
	ExpiresIn     int64
}
