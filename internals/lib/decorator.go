package lib

import (
	"context"
	"log"
	"time"

	"github.com/himanshukumar42/soundauth/internals/models"
	"github.com/himanshukumar42/soundauth/internals/provider"
)

type LoggingDecorator struct {
	next provider.AuthenticationProvider
}

func NewLoggingDecorator(next provider.AuthenticationProvider) *LoggingDecorator {
	return &LoggingDecorator{
		next: next,
	}
}

func (ld *LoggingDecorator) Name() string {
	return ld.next.Name()
}

func (ld *LoggingDecorator) Log(ctx context.Context, req models.AuthRequest) (*models.AuthResponse, error) {
	start := time.Now()

	log.Printf("[LOG] Authentication started Provider=%s Tenant=%s", req.Provider, req.TenantID)

	resp, err := ld.next.Authenticate(ctx, req)

	log.Printf("[LOG] Authentication finished provider=%s Duration=%s", req.Provider, time.Since(start))

	return resp, err
}

type MetricsDecorator struct {
	next provider.AuthenticationProvider
}

func NewMetricsDecorator(next provider.AuthenticationProvider) *MetricsDecorator {
	return &MetricsDecorator{
		next: next,
	}
}

func (d *MetricsDecorator) Name() string {
	return d.next.Name()
}

func (md *MetricsDecorator) Metric(ctx context.Context, req models.AuthRequest) (*models.AuthResponse, error) {
	start := time.Now()

	resp, err := md.next.Authenticate(ctx, req)

	log.Printf("[METRICS] provider=%s latency=%s success=%v", req.Provider, time.Since(start), err == nil)

	return resp, err
}

type AuditDecorator struct {
	next provider.AuthenticationProvider
}

func NewAuditDecorator(next provider.AuthenticationProvider) *AuditDecorator {
	return &AuditDecorator{
		next: next,
	}
}

func (ad *AuditDecorator) Name() string {
	return ad.next.Name()
}

func (ad *AuditDecorator) Audit(ctx context.Context, req models.AuthRequest) (*models.AuthResponse, error) {
	log.Printf("[AUDIT] login attempt tenat=%s provider=%s device=%s", req.TenantID, req.Provider, req.DeviceID)

	return ad.next.Authenticate(ctx, req)
}
