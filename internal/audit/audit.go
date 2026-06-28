package audit

import (
	"context"
	"time"
)

type AuditEvent struct {
	ID        string                 `json:"id"`
	TenantID  string                 `json:"tenant_id"`
	UserID    string                 `json:"user_id"`
	EventType string                 `json:"event_type"`
	Timestamp time.Time              `json:"timestamp"`
	Actor     string                 `json:"actor"` // Email or System
	Success   bool                   `json:"success"`
	IPAddress string                 `json:"ip_address,omitempty"`
	UserAgent string                 `json:"user_agent,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

type QueryFilter struct {
	EventType string
	Success   *bool
	Limit     int
	Offset    int
}

type AuditLogRepository interface {
	LogEvent(ctx context.Context, event *AuditEvent) error
	QueryEvents(ctx context.Context, tenantID string, filter QueryFilter) ([]AuditEvent, error)
}
