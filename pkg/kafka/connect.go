package kafka

import (
	"github.com/IBM/sarama"
	"messageprocessor/internal/config"
)

func NewSyncProducer(cfg config.ConfigInterface) (sarama.SyncProducer, error) {
	cfgSarama := sarama.NewConfig()

	cfgSarama.Producer.Partitioner = sarama.NewRandomPartitioner
	cfgSarama.Producer.RequiredAcks = sarama.WaitForAll
	cfgSarama.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(cfg.GetKafka().Brockers, cfgSarama)
	return producer, err
}

func NewAsyncProducer(cfg config.ConfigInterface) (sarama.AsyncProducer, error) {
	cfgSarama := sarama.NewConfig()

	cfgSarama.Producer.Partitioner = sarama.NewRandomPartitioner
	cfgSarama.Producer.RequiredAcks = sarama.WaitForAll
	cfgSarama.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer(cfg.GetKafka().Brockers, cfgSarama)
	return producer, err
}

func PrepareMessage(topic string, message []byte) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
		Value:     sarama.ByteEncoder(message),
	}
	return msg
}
