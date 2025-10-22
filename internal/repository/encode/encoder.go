package encode

import (
	"encoding/json"
)

type JSONEncode struct{}

func New() *JSONEncode {
	return &JSONEncode{}
}

func (e *JSONEncode) Decode(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (e *JSONEncode) Encode(v any) ([]byte, error) {
	return json.Marshal(v)
}
