package kafka

import (
	"github.com/Shopify/sarama"
)

func NewProducer(broker string) sarama.SyncProducer {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{broker}, config)
	if err != nil {
		panic(err)
	}
	return producer
}

func PrepareMessage(message string, topic string) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(message)}
	return msg
}
