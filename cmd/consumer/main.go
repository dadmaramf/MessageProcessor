package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"messageprocessor/internal/config"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Message struct {
	ID      int    `json:"id" db:"id"`
	Content string `json:"content" db:"content"`
}

type Consumer struct{}

func (c *Consumer) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var msg Message
		if err := json.Unmarshal(message.Value, &msg); err != nil {
			fmt.Printf("could not unmarshal message: %s", err)
			return err
		}
		fmt.Println(string(msg.Content))
		session.MarkMessage(message, "")
	}
	return nil
}

func subscribe(ctx context.Context, topic string, consumerGroup sarama.ConsumerGroup, wg *sync.WaitGroup) error {
	consumer := Consumer{}
	go func() {
		defer wg.Done()
		if err := consumerGroup.Consume(ctx, []string{topic}, &consumer); err != nil {
			fmt.Printf("Error consemer: %v", err)
		}

		if ctx.Err() != nil {
			return
		}
	}()

	return nil
}

func StartConsumer(ctx context.Context, wg *sync.WaitGroup) error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}
	configSarama := sarama.NewConfig()

	configSarama.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	configSarama.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumerGroup, err := sarama.NewConsumerGroup(cfg.Kf.Brockers, "analytic", configSarama)

	if err != nil {
		return err
	}
	wg.Add(2)
	err = subscribe(ctx, "messages", consumerGroup, wg)
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Second)
	return subscribe(ctx, "messages", consumerGroup, wg)
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()
	var wg sync.WaitGroup
	err := StartConsumer(ctx, &wg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
	}

	wg.Wait()
}
