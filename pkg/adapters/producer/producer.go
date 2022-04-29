package producer

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	PRODUCER_KAFKA = "KAFKA"
)

type ProducerConfig struct {
	Username string
	Password string
	Path     string
	Scheme   string
	Host     string
	TLS      bool
}

type Producer interface {
	SendMessage(message []byte, destination string) error
}

func NewProducer(kind string, host string, tls bool) (Producer, error) {
	var prod Producer

	config, err := getConfig(host, tls)
	if err != nil {
		return prod, err
	}

	switch strings.ToUpper(kind) {
	case PRODUCER_KAFKA:
		prod, err = NewKafkaProducer(config)
	default:
		err = fmt.Errorf("No valid producer type passed")
	}
	return prod, err

}

func getConfig(host string, tls bool) (*ProducerConfig, error) {
	urlInfo, err := url.Parse(host)
	if err != nil {
		return &ProducerConfig{}, err
	}

	password, _ := urlInfo.User.Password()

	return &ProducerConfig{
		Username: urlInfo.User.Username(),
		Password: password,
		Path:     urlInfo.Path,
		Scheme:   urlInfo.Scheme,
		Host:     urlInfo.Host,
		TLS:      tls,
	}, nil
}

func IsValidProducerType(kind string) bool {
	kind = strings.ToUpper(kind)
	return kind == PRODUCER_KAFKA
}
