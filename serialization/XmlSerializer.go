package serialization

import (
	"encoding/xml"
)

type XmlSerializer struct{}

// Marshal implements ISerializer.
func (j *XmlSerializer) Marshal(v any) ([]byte, error) {
	return xml.Marshal(v)
}

// Unmarshal implements ISerializer.
func (j *XmlSerializer) Unmarshal(data []byte, v any) error {
	return xml.Unmarshal(data, v)
}

func NewXmlSerializer() ISerializer {
	return &XmlSerializer{}
}
