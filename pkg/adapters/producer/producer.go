package producer

import (
	"fmt"
	"strings"
)

const (
	PRODUCER_KAFKA = "KAFKA"
)

type Producer interface {
	SendMessage(message []byte, destination string) error
}

func NewProducer(kind string, host string, encoding string) (Producer, error) {
	var prod Producer
	var err error
	switch strings.ToUpper(kind) {
	case PRODUCER_KAFKA:
		prod, err = NewKafkaProducer(host, encoding)
	default:
		err = fmt.Errorf("No valid producer type passed")
	}
	return prod, err

}

func IsValidProducerType(kind string) bool {
	kind = strings.ToUpper(kind)
	return kind == PRODUCER_KAFKA
}
