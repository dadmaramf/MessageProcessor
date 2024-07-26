package storage

import (
	"context"
	"messageprocessor/internal/model"
)

type Storage interface {
	GetNewOutbox(ctx context.Context) (*model.Message, error)
	SetDown(id int) error
	PostMessage(msg string) (err error)
	GetDownMessages() ([]model.Message, error)
}
