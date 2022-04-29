package producer

import (
	"github.com/Shopify/sarama"
	"github.com/tianet/ditto/pkg/adapters/encoder"
)

type KafkaProducer struct {
	syncProducer sarama.SyncProducer
	encoding     string
}

func NewKafkaProducer(broker string, enc string) (KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{broker}, config)
	if err != nil {
		return KafkaProducer{}, err
	}
	return KafkaProducer{syncProducer: producer, encoding: enc}, nil
}

func (p KafkaProducer) SendMessage(message []byte, destination string) error {
	var value sarama.Encoder
	switch p.encoding {
	case encoder.JSON:
		value = sarama.StringEncoder(message)
	case encoder.AVRO:
		value = sarama.ByteEncoder(message)
	}

	msg := &sarama.ProducerMessage{Topic: destination, Value: value}
	_, _, err := p.syncProducer.SendMessage(msg)
	if err != nil {
		return err
	}
	return nil
}
