package services

import (
	"context"
	"messageprocessor/internal/model"
	"time"
)

//go:generate mockgen -source=service.go -destination=mock/service.go
type Service interface {
	StartProcessingMessage(ctx context.Context, handlePeriod time.Duration)
	SaveMessage(msg string) error
	SentMessages() ([]model.MessageState, error)
}
