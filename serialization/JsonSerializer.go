package serialization

import "encoding/json"

type JsonSerializer struct{}

// Marshal implements ISerializer.
func (j *JsonSerializer) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal implements ISerializer.
func (j *JsonSerializer) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func NewJSONSerializer() ISerializer {
	return &JsonSerializer{}
}
