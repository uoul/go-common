package serialization

import "gopkg.in/yaml.v3"

type YamlSerializer struct{}

// Marshal implements ISerializer.
func (j *YamlSerializer) Marshal(v any) ([]byte, error) {
	return yaml.Marshal(v)
}

// Unmarshal implements ISerializer.
func (j *YamlSerializer) Unmarshal(data []byte, v any) error {
	return yaml.Unmarshal(data, v)
}

func NewYamlSerializer() ISerializer {
	return &YamlSerializer{}
}
