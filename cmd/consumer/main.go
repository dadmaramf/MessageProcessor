package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

type Message struct {
	ID        int
	Content   []byte
	CreatedAt time.Time
}

type Consumer struct{}

func (c *Consumer) Setup() {
}

func (c *Consumer) c() {

}

func (c *Consumer) ConsumeClain(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) {
	for message := range claim.Messages() {
		var msg Message
		if err := json.Unmarshal(message.Value, &msg); err != nil {
			fmt.Printf("could not unmarshal message: %s", err)
		}
		fmt.Println(string(msg.Content))
		session.MarkMessage(message, "")
	}
}
func main() {

}
