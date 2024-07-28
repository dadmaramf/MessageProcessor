package services

import (
	"context"
	"github.com/IBM/sarama"
	"log/slog"
	"messageprocessor/internal/model"
	messagereader "messageprocessor/internal/services/message_reader"
	messagesender "messageprocessor/internal/services/message_sender"
	"messageprocessor/internal/storage"
	"time"
)

//go:generate mockgen -source=service.go -destination=mock/service.go

type SyncProducer interface {
	sarama.SyncProducer
}

type ConsumerGroup interface {
	sarama.ConsumerGroup
}

type Service interface {
	messagesender.ServiceSender
	messagereader.ServiceReader
}

type Services struct {
	s messagesender.ServiceSender
	r messagereader.ServiceReader
}

func NewServices(storage storage.Storage, cons sarama.ConsumerGroup, prod sarama.SyncProducer, log *slog.Logger) *Services {
	return &Services{
		s: messagesender.New(storage, prod, log),
		r: messagereader.New(cons, storage, log),
	}
}

func (s *Services) SaveMessage(msg string) error {
	err := s.s.SaveMessage(msg)
	return err
}

func (s *Services) SentMessages() ([]model.Message, error) {
	msgs, err := s.s.SentMessages()
	return msgs, err
}

func (s *Services) StartConsumerProcessingMessage(ctx context.Context) {
	s.r.StartConsumerProcessingMessage(ctx)
}

func (s *Services) StartProcessingMessage(ctx context.Context, handlePeriod time.Duration) {
	s.s.StartProcessingMessage(ctx, handlePeriod)
}
