package server

import (
	"fmt"
	"strings"

	"github.com/tianet/ditto/pkg/application/schema"
)

const (
	SERVER_REST = "REST"
)

type Server interface {
	AddEndpoint(endpoint string, fields *[]schema.Field)
	ListenAndServe()
}

func NewServer(kind string, port int) (Server, error) {
	switch strings.ToUpper(kind) {
	case SERVER_REST:
		return NewRestServer(port)
	default:
		return nil, fmt.Errorf("Server type %s not supported", kind)
	}
}

func IsValidServerType(kind string) bool {
	kind = strings.ToUpper(kind)
	return kind == SERVER_REST
}
