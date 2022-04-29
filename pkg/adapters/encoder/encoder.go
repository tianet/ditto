package encoder

const (
	JSON = "JSON"
	AVRO = "AVRO"
)

type Encoder interface {
	Marshal(map[string]interface{}) ([]byte, error)
}

func NewEncoder(_type string, schemaPath string) (Encoder, error) {
	var enc Encoder
	var err error

	switch _type {
	case JSON:
		enc, err = NewJsonEncoder()
	case AVRO:
		enc, err = NewAvroEncoder(schemaPath)
	}
	return enc, err
}
