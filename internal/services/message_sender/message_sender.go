package messagesender

import (
	"context"
	"log/slog"
	"messageprocessor/internal/model"
	"messageprocessor/internal/services"
	"messageprocessor/internal/storage"
	"time"
)

var _ services.Service = (*Sender)(nil)

type Sender struct {
	storage storage.Storage
	log     *slog.Logger
}

func New(storage storage.Storage, log *slog.Logger) *Sender {
	return &Sender{
		storage: storage,
		log:     log,
	}
}

func (s *Sender) SaveMessage(msg string) error {
	return s.storage.PostMessage(msg)
}

func (s *Sender) SentMessages() ([]model.Message, error) {
	msgs, err := s.storage.GetDownMessages()
	return msgs, err
}

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
			if err := s.storage.SetDown(message.ID); err != nil {
				log.Error("failed to det message done")
			}
		}
	}()
}

func (s *Sender) sendMessage(msg model.Message) {
	// const op = "services.message-sender.SendMessage"
	// log := s.log.With(slog.String("op", op))

}
