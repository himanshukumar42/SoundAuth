package models

type ProviderResponse struct {
	UserID        string
	Authenticated bool
	Metadata      map[string]string
}
