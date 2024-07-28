package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"messageprocessor/internal/config"
	"messageprocessor/pkg/kafka"
	"os"
	"sync"

	"github.com/IBM/sarama"
)

type Message struct {
	ID      int    `json:"id" db:"id"`
	Content string `json:"content" db:"content"`
}

type Consumer struct {
	logger   *slog.Logger
	producer sarama.SyncProducer
}

func (c *Consumer) Setup(session sarama.ConsumerGroupSession) error {
	c.logger.Info("Consumer group session setup")
	return nil
}

func (c *Consumer) Cleanup(session sarama.ConsumerGroupSession) error {
	c.logger.Info("Consumer group session cleanup")
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var msg Message
		if err := json.Unmarshal(message.Value, &msg); err != nil {
			c.logger.Error("Could not unmarshal message", slog.Any("error", err))
			return err
		}
		c.logger.Info("Message received", slog.String("content", msg.Content))
		session.MarkMessage(message, "")

		msg.Content += ". Processed!"

		msgByte, err := json.Marshal(msg)

		if err != nil {
			c.logger.Error("Could not unmarshal message", slog.Any("error", err))
			return err
		}

		resp := kafka.PrepareMessage("processedMessage", msgByte)

		_, _, err = c.producer.SendMessage(resp)

		if err != nil {
			c.logger.Error("Could not send message", slog.Any("error", err))
			return err
		}

	}
	return nil
}

func subscribe(
	ctx context.Context,
	topic string,
	consumerGroup sarama.ConsumerGroup,
	wg *sync.WaitGroup,
	logger *slog.Logger,
	producer sarama.SyncProducer,
) error {

	consumer := &Consumer{logger: logger, producer: producer}
	go func() {
		defer wg.Done()
		for {
			if err := consumerGroup.Consume(ctx, []string{topic}, consumer); err != nil {
				logger.Error("Error consuming", slog.Any("error", err))
			}

			if ctx.Err() != nil {
				return
			}
		}
	}()
	return nil
}

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	producer, err := kafka.NewSyncProducer(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer producer.Close()

	consumer, err := kafka.NewConsumerGroup(cfg, "analytics")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer consumer.Close()

	ctx := context.Background()
	wg := &sync.WaitGroup{}
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	wg.Add(1)
	subscribe(ctx, "messages", consumer, wg, log, producer)
	wg.Wait()
}
