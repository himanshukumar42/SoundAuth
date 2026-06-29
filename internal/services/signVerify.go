package services

import (
	"context"
	"log"
	"time"

	"github.com/himanshukumar42/soundauth/internal/models"
	worker "github.com/himanshukumar42/soundauth/internal/workerpool"
)

type SignatureVerifier struct {
	tenantService *TenantService
	pool          *worker.VerificationPool
}

func NewSignatureVerifier(tenantService *TenantService, pool *worker.VerificationPool) *SignatureVerifier {
	return &SignatureVerifier{
		tenantService: tenantService,
		pool:          pool,
	}
}

func (s *SignatureVerifier) Verify(ctx context.Context, req models.AuthRequest) error {
	signingKey, err := s.tenantService.GetSigninKey(ctx, req.TenantID)
	if err != nil {
		return err
	}

	log.Printf("[Verifier] Loaded signing key %s", signingKey)

	jobs := []worker.VerificationJob{
		{
			Name: "Challenge",
			Run: func(ctx context.Context) error {
				time.Sleep(100 * time.Millisecond)
				log.Println("Challenge Verified")
				return nil
			},
		},
		{
			Name: "Device",
			Run: func(ctx context.Context) error {
				time.Sleep(100 * time.Millisecond)
				log.Printf("Device %s verified", req.DeviceID)
				return nil
			},
		},
		{
			Name: "Public Key",
			Run: func(ctx context.Context) error {
				time.Sleep(200 * time.Millisecond)
				log.Println("Public Key Verified")
				return nil
			},
		},
		{
			Name: "Origin",
			Run: func(ctx context.Context) error {
				time.Sleep(120 * time.Millisecond)
				log.Println("Origin Verified")
				return nil
			},
		},
	}
	return s.pool.Verify(ctx, jobs)
}
