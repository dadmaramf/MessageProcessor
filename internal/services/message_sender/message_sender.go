package messagesender

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log/slog"
	"messageprocessor/internal/model"
	"messageprocessor/internal/services"
	"messageprocessor/internal/storage"
	"messageprocessor/pkg/kafka"
	"time"
)

var _ services.Service = (*Sender)(nil)

type Sender struct {
	storage storage.Storage
	prod    sarama.SyncProducer
	log     *slog.Logger
}

func New(storage storage.Storage, prod sarama.SyncProducer, log *slog.Logger) *Sender {
	return &Sender{
		storage: storage,
		prod:    prod,
		log:     log,
	}
}

func (s *Sender) SaveMessage(msg string) error {
	return s.storage.PostMessage(msg)
}

func (s *Sender) SentMessages() ([]model.MessageState, error) {
	msgs, err := s.storage.GetDownMessages()
	return msgs, err
}

// StartProcessingMessage starts the process of handling messages at regular intervals.
func (s *Sender) StartProcessingMessage(ctx context.Context, handlePeriod time.Duration) {
	const op = "services.message-sender.StartProcessingMessage"

	log := s.log.With(slog.String("op", op))

	tiker := time.NewTicker(handlePeriod)

	go func() {
		for {
			select {
			case <-ctx.Done():
				s.log.Info("stopping send processing")
				tiker.Stop()
				return
			case <-tiker.C:
			}
			message, err := s.storage.GetNewOutbox(ctx)
			if err != nil {
				log.Error("failed to get new message", "error", slog.StringValue(err.Error()))
				continue
			}

			if message == nil {
				log.Debug("no new message")
				continue
			}
			msgByte, err := s.encode(message)
			if err != nil {
				log.Error("failed to message convert", "error", slog.StringValue(err.Error()))
				continue
			}
			if err := s.sendMessage(msgByte); err != nil {
				log.Error("failed to send message", "error", slog.StringValue(err.Error()))
				continue
			}
			if err := s.storage.SetDown(message.ID); err != nil {
				log.Error("failed to det message done", "error", err)
			}
			log.Info("message send")
		}
	}()
}

// sendMessage sends a message to Kafka.
func (s *Sender) sendMessage(msg []byte) error {

	m := kafka.PrepareMessage("messages", msg)
	_, _, err := s.prod.SendMessage(m)
	return err
}

func (s *Sender) encode(msg *model.Message) ([]byte, error) {
	msgByte, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("encode json: %w", err)
	}
	return msgByte, nil
}
