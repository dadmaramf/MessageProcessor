package messagesender

import (
	"context"
	"messageprocessor/internal/model"
	"time"
)

//go:generate mockgen -source=service.go -destination=../mock/servicesender.go
type ServiceSender interface {
	StartProcessingMessage(ctx context.Context, handlePeriod time.Duration)
	SaveMessage(msg string) error
	SentMessages() ([]model.Message, error)
}
