package tools

import (
	"path/filepath"
	"strings"
)

func GetTopic(path string) string {
	return strings.Split(filepath.Base(path), ".")[0]
}

func GetEndpoint(path string) string {
	return "/" + GetTopic(path)
}
