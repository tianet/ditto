package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/tianet/ditto/pkg/schema"
)

type Endpoint struct {
	Fields *[]schema.Field
}

func (s Endpoint) Handler(w http.ResponseWriter, req *http.Request) {
	message, err := schema.GenerateMessage(s.Fields)

	output, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, string(output))
}

func GetEndpoint(path string) string {
	return "/" + strings.Split(filepath.Base(path), ".")[0]
}
