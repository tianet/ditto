package encoder

import (
	"os"
	"path/filepath"

	"github.com/linkedin/goavro/v2"
)

type AvroEncoder struct {
	codec *goavro.Codec
}

func NewAvroEncoder(path string) (AvroEncoder, error) {
	schema, err := getAvroSchemaFromFile(path)
	if err != nil {
		return AvroEncoder{}, err
	}

	codec, err := goavro.NewCodec(schema)
	if err != nil {
		return AvroEncoder{}, err
	}
	return AvroEncoder{
		codec: codec,
	}, nil
}

func getAvroSchemaFromFile(path string) (string, error) {
	data, err := os.ReadFile(filepath.Join(path))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (a AvroEncoder) Marshal(message map[string]interface{}) ([]byte, error) {
	return a.codec.BinaryFromNative(nil, message)
}
