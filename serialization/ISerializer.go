package serialization

type ISerializer interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}
