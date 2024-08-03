package config

type IConfigProvider interface {
	SetValue(key string, value string)

	// This methods should return the current config value, if exist, otherwise the given default value
	StringOrDefault(key string, defaultValue string) string
	IntOrDefault(key string, defaultValue int) int
	Int8OrDefault(key string, defaultValue int8) int8
	Int16OrDefault(key string, defaultValue int16) int16
	Int32OrDefault(key string, defaultValue int32) int32
	Int64OrDefault(key string, defaultValue int64) int64
	UIntOrDefault(key string, defaultValue uint) uint
	UInt8OrDefault(key string, defaultValue uint8) uint8
	UInt16OrDefault(key string, defaultValue uint16) uint16
	UInt32OrDefault(key string, defaultValue uint32) uint32
	UInt64OrDefault(key string, defaultValue uint64) uint64

	// This methods should return the config value, or error
	String(key string) (string, error)
	Int(key string) (int, error)
	Int8(key string) (int8, error)
	Int16(key string) (int16, error)
	Int32(key string) (int32, error)
	Int64(key string) (int64, error)
	UInt(key string) (uint, error)
	UInt8(key string) (uint8, error)
	UInt16(key string) (uint16, error)
	UInt32(key string) (uint32, error)
	UInt64(key string) (uint64, error)
}
