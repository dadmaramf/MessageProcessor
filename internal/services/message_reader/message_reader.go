package messagereader

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"log/slog"
	"messageprocessor/internal/storage"
)

var _ ServiceReader = (*Reader)(nil)

type Reader struct {
	cons    sarama.ConsumerGroup
	storage storage.Storage
	log     *slog.Logger
}

func New(cons sarama.ConsumerGroup, storage storage.Storage, log *slog.Logger) *Reader {
	return &Reader{
		cons:    cons,
		storage: storage,
		log:     log,
	}
}

func (r *Reader) StartConsumerProcessingMessage(ctx context.Context) {
	const op = "services.message-sender.StartConsumingProcessedMessages"
	log := r.log.With(slog.String("op", op))
	consumer := NewConsumer(r.log, r.storage)

	go func() {
		for {
			if err := r.cons.Consume(ctx, []string{"processedMessage"}, consumer); err != nil {
				log.Error("Error consuming", "error", err)
			}

			if ctx.Err() != nil {
				return
			}
		}
	}()

}

type Consumer struct {
	log  *slog.Logger
	stor storage.Storage
}

func NewConsumer(log *slog.Logger, stor storage.Storage) *Consumer {
	return &Consumer{log: log, stor: stor}
}

func (c *Consumer) Setup(session sarama.ConsumerGroupSession) error {
	c.log.Info("Consumer group session setup")
	return nil
}

func (c *Consumer) Cleanup(session sarama.ConsumerGroupSession) error {
	c.log.Info("Consumer group session cleanup")
	return nil
}

type Message struct {
	ID      int    `json:"id" db:"id"`
	Content string `json:"content" db:"content"`
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var msg Message
		if err := json.Unmarshal(message.Value, &msg); err != nil {
			c.log.Error("Could not unmarshal message", "error", slog.StringValue(err.Error()))
			continue
		}
		if err := c.processMessage(msg.ID, msg.Content); err != nil {
			c.log.Error("Failed to process message", "error", slog.StringValue(err.Error()))
			continue
		}
		session.MarkMessage(message, "")
	}
	return nil
}

func (c *Consumer) processMessage(id int, msg string) error {
	err := c.stor.AddProcessedMessage(id, msg)
	if err != nil {
		return err
	}
	return nil
}
