package encoder

import "encoding/json"

type JsonEncoder struct{}

func NewJsonEncoder() (JsonEncoder, error) {
	return JsonEncoder{}, nil
}

func (j JsonEncoder) Marshal(message map[string]interface{}) ([]byte, error) {
	return json.Marshal(message)
}
