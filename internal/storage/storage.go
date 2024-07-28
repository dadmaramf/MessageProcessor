package storage

import (
	"context"
	"messageprocessor/internal/model"
)

//go:generate mockgen -source=storage.go -destination=mock/storage.go
type Storage interface {
	GetNewOutbox(ctx context.Context) (*model.Message, error)
	SetDown(id int) error
	PostMessage(msg string) (err error)
	GetDownMessages() ([]model.Message, error)
	AddProcessedMessage(id int, msg string) error
}
