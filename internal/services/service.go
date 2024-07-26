package services

import (
	"context"
	"messageprocessor/internal/model"
	"time"
)

type Service interface {
	StartProcessingMessage(ctx context.Context, handlePeriod time.Duration)
	SaveMessage(msg string) error
	SentMessages() ([]model.Message, error)
}
