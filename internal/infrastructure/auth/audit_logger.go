package auth

import (
	"clean-architecture-api/internal/domain/repositories"
	"clean-architecture-api/pkg/logger"
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditLogEntry struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Action    string    `json:"action" gorm:"not null"`
	Resource  string    `json:"resource" gorm:"not null"`
	EntityID  uuid.UUID `json:"entity_id" gorm:"type:uuid"`
	Timestamp time.Time `json:"timestamp" gorm:"not null"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
}

func (ale *AuditLogEntry) BeforeCreate(_ *gorm.DB) error {
	if ale.ID == uuid.Nil {
		ale.ID = uuid.New()
	}
	return nil
}

type AuditLoggerImpl struct {
	logger logger.Logger
}

func NewAuditLogger(logger logger.Logger) repositories.AuditLogger {
	return &AuditLoggerImpl{
		logger: logger,
	}
}

func (a *AuditLoggerImpl) LogAccess(_ context.Context, userID uuid.UUID, action, resource string, entityID uuid.UUID) error {
	entry := AuditLogEntry{
		ID:        uuid.New(),
		UserID:    userID,
		Action:    action,
		Resource:  resource,
		EntityID:  entityID,
		Timestamp: time.Now(),
	}

	a.logger.WithField("user_id", entry.UserID).
		WithField("action", entry.Action).
		WithField("resource", entry.Resource).
		WithField("entity_id", entry.EntityID).
		WithField("timestamp", entry.Timestamp).
		Info("Audit log entry")

	return nil
}

func (a *AuditLoggerImpl) LogDataAccess(ctx context.Context, userID uuid.UUID, action, resource string, data interface{}) error {
	entry := AuditLogEntry{
		ID:        uuid.New(),
		UserID:    userID,
		Action:    action,
		Resource:  resource,
		Timestamp: time.Now(),
	}

	a.logger.WithField("user_id", entry.UserID).
		WithField("action", entry.Action).
		WithField("resource", entry.Resource).
		WithField("data", data).
		WithField("timestamp", entry.Timestamp).
		Info("Data access audit log")

	return nil
}
