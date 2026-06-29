package models

type Provider string

const (
	ProviderPasskey Provider = "Passkey"
	ProviderGoogle  Provider = "GoogleOAuth"
	ProviderGithub  Provider = "GithubOAuth"
	ProviderMagic   Provider = "MagicLink"
	ProviderSAML    Provider = "SAML"
)

type ProviderResponse struct {
	UserID        string
	Authenticated bool
	Metadata      map[string]string
}
