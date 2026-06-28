package vault

import (
	"context"
	"fmt"
	"sync"
)

type VaultClient struct {
	mu      sync.RWMutex
	Secrets map[string]string
}

func NewVaultClient() *VaultClient {
	return &VaultClient{
		Secrets: map[string]string{
			"google/signing-key":  "GOOGLE_SECRET_KEY",
			"github/signing-key":  "GITHUB_SECRET_KEY",
			"saml/signing-key":    "SAML_SECRET_KEY",
			"magic/signing-key":   "MAGIC_SECRET_KEY",
			"passkey/signing-key": "PASSKEY_SECRET_KEY",
		},
	}
}

func (vc *VaultClient) GetSecret(ctx context.Context, path string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	vc.mu.RLock()
	defer vc.mu.RUnlock()

	secret, ok := vc.Secrets[path]
	if !ok {
		return "", fmt.Errorf("secret %s not found", path)
	}
	return secret, nil
}

