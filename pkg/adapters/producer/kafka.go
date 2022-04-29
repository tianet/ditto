package producer

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type KafkaProducer struct {
	config *ProducerConfig
	writer *kafka.Writer
}

const (
	PLAIN       = "PLAIN"
	SCRAMSHA256 = "SCRAM-SHA-256"
	SCRAMSHA512 = "SCRAM-SHA-512"
)

func NewKafkaProducer(config *ProducerConfig) (KafkaProducer, error) {
	transport := kafka.DefaultTransport
	tls := &tls.Config{InsecureSkipVerify: true}

	if config.TLS {
		transport = &kafka.Transport{TLS: tls}
	}
	var err error
	if config.Username != "" && config.Password != "" {

		var mechanism sasl.Mechanism
		switch strings.ToUpper(config.Scheme) {
		case PLAIN:
			mechanism = plain.Mechanism{Username: config.Username, Password: config.Password}
		case SCRAMSHA256:
			mechanism, err = scram.Mechanism(scram.SHA256, config.Username, config.Password)
			if err != nil {
				return KafkaProducer{}, err
			}
		case SCRAMSHA512:
			mechanism, err = scram.Mechanism(scram.SHA512, config.Username, config.Password)
			if err != nil {
				return KafkaProducer{}, err
			}
		default:
			return KafkaProducer{}, fmt.Errorf("Kafka mechanism %s not supported", config.Scheme)
		}

		if config.TLS {
			transport = &kafka.Transport{
				SASL: mechanism,
				TLS:  tls,
			}
		} else {
			transport = &kafka.Transport{
				SASL: mechanism,
			}

		}
	}

	if err != nil {
		return KafkaProducer{}, err
	}

	return KafkaProducer{
		config: config,
		writer: &kafka.Writer{
			Addr:      kafka.TCP(config.Host),
			Transport: transport,
		},
	}, nil
}

func (p KafkaProducer) SendMessage(message []byte, destination string) error {
	msg := kafka.Message{
		Topic: destination,
		Value: message,
	}
	err := p.writer.WriteMessages(context.Background(), msg)
	if err != nil {
		return err
	}
	return nil
}
