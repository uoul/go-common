package config

import (
	"math/bits"
	"os"
	"strconv"
	"strings"
)

type EnvVarProvider struct {
	config map[string]string
}

// Int implements IConfigProvider.
func (e *EnvVarProvider) Int(key string) (int, error) {
	if configValue, keyExists := e.config[key]; keyExists {
		retVal, err := strconv.ParseInt(configValue, 10, bits.UintSize)
		if err != nil {
			return 0, NewConfigError("failed to parse %s as int of parameter with key %s", configValue, key)
		}
		return int(retVal), nil
	}
	return 0, NewConfigError("key %s does not exist", key)
}

// Int32 implements IConfigProvider.
func (e *EnvVarProvider) Int32(key string) (int32, error) {
	if configValue, keyExists := e.config[key]; keyExists {
		retVal, err := strconv.ParseInt(configValue, 10, 32)
		if err != nil {
			return 0, NewConfigError("failed to parse %s as int32 of parameter with key %s", configValue, key)
		}
		return int32(retVal), nil
	}
	return 0, NewConfigError("key %s does not exist", key)
}

// Int64 implements IConfigProvider.
func (e *EnvVarProvider) Int64(key string) (int64, error) {
	if configValue, keyExists := e.config[key]; keyExists {
		retVal, err := strconv.ParseInt(configValue, 10, 64)
		if err != nil {
			return 0, NewConfigError("failed to parse %s as int64 of parameter with key %s", configValue, key)
		}
		return int64(retVal), nil
	}
	return 0, NewConfigError("key %s does not exist", key)
}

// Int8 implements IConfigProvider.
func (e *EnvVarProvider) Int8(key string) (int8, error) {
	if configValue, keyExists := e.config[key]; keyExists {
		retVal, err := strconv.ParseInt(configValue, 10, 8)
		if err != nil {
			return 0, NewConfigError("failed to parse %s as int8 of parameter with key %s", configValue, key)
		}
		return int8(retVal), nil
	}
	return 0, NewConfigError("key %s does not exist", key)
}

// Int16 implements IConfigProvider.
func (e *EnvVarProvider) Int16(key string) (int16, error) {
	if configValue, keyExists := e.config[key]; keyExists {
		retVal, err := strconv.ParseInt(configValue, 10, 16)
		if err != nil {
			return 0, NewConfigError("failed to parse %s as int16 of parameter with key %s", configValue, key)
		}
		return int16(retVal), nil
	}
	return 0, NewConfigError("key %s does not exist", key)
}

// String implements IConfigProvider.
func (e *EnvVarProvider) String(key string) (string, error) {
	if configValue, keyExists := e.config[key]; keyExists {
		return configValue, nil
	}
	return "", NewConfigError("key %s does not exist", key)
}

// UInt implements IConfigProvider.
func (e *EnvVarProvider) UInt(key string) (uint, error) {
	if configValue, keyExists := e.config[key]; keyExists {
		retVal, err := strconv.ParseUint(configValue, 10, bits.UintSize)
		if err != nil {
			return 0, NewConfigError("failed to parse %s as uint of parameter with key %s", configValue, key)
		}
		return uint(retVal), nil
	}
	return 0, NewConfigError("key %s does not exist", key)
}

// UInt32 implements IConfigProvider.
func (e *EnvVarProvider) UInt32(key string) (uint32, error) {
	if configValue, keyExists := e.config[key]; keyExists {
		retVal, err := strconv.ParseUint(configValue, 10, 32)
		if err != nil {
			return 0, NewConfigError("failed to parse %s as uint32 of parameter with key %s", configValue, key)
		}
		return uint32(retVal), nil
	}
	return 0, NewConfigError("key %s does not exist", key)
}

// UInt64 implements IConfigProvider.
func (e *EnvVarProvider) UInt64(key string) (uint64, error) {
	if configValue, keyExists := e.config[key]; keyExists {
		retVal, err := strconv.ParseUint(configValue, 10, 64)
		if err != nil {
			return 0, NewConfigError("failed to parse %s as uint64 of parameter with key %s", configValue, key)
		}
		return uint64(retVal), nil
	}
	return 0, NewConfigError("key %s does not exist", key)
}

// UInt8 implements IConfigProvider.
func (e *EnvVarProvider) UInt8(key string) (uint8, error) {
	if configValue, keyExists := e.config[key]; keyExists {
		retVal, err := strconv.ParseUint(configValue, 10, 8)
		if err != nil {
			return 0, NewConfigError("failed to parse %s as uint8 of parameter with key %s", configValue, key)
		}
		return uint8(retVal), nil
	}
	return 0, NewConfigError("key %s does not exist", key)
}

// UInt16 implements IConfigProvider.
func (e *EnvVarProvider) UInt16(key string) (uint16, error) {
	if configValue, keyExists := e.config[key]; keyExists {
		retVal, err := strconv.ParseUint(configValue, 10, 16)
		if err != nil {
			return 0, NewConfigError("failed to parse %s as uint16 of parameter with key %s", configValue, key)
		}
		return uint16(retVal), nil
	}
	return 0, NewConfigError("key %s does not exist", key)
}

// GetInt implements IConfigProvider.
func (e *EnvVarProvider) IntOrDefault(key string, defaultValue int) int {
	if configValue, keyExists := e.config[key]; keyExists {
		retVal, err := strconv.ParseInt(configValue, 10, bits.UintSize)
		if err != nil {
			return defaultValue
		}
		return int(retVal)
	}
	return defaultValue
}

// GetInt16 implements IConfigProvider.
func (e *EnvVarProvider) Int16OrDefault(key string, defaultValue int16) int16 {
	if configValue, keyExists := e.config[key]; keyExists {
		retVal, err := strconv.ParseInt(configValue, 10, 16)
		if err != nil {
			return defaultValue
		}
		return int16(retVal)
	}
	return defaultValue
}

// GetInt32 implements IConfigProvider.
func (e *EnvVarProvider) Int32OrDefault(key string, defaultValue int32) int32 {
	if configValue, keyExists := e.config[key]; keyExists {
		retVal, err := strconv.ParseInt(configValue, 10, 32)
		if err != nil {
			return defaultValue
		}
		return int32(retVal)
	}
	return defaultValue
}

// GetInt64 implements IConfigProvider.
func (e *EnvVarProvider) Int64OrDefault(key string, defaultValue int64) int64 {
	if configValue, keyExists := e.config[key]; keyExists {
		retVal, err := strconv.ParseInt(configValue, 10, 64)
		if err != nil {
			return defaultValue
		}
		return int64(retVal)
	}
	return defaultValue
}

// GetInt8 implements IConfigProvider.
func (e *EnvVarProvider) Int8OrDefault(key string, defaultValue int8) int8 {
	if configValue, keyExists := e.config[key]; keyExists {
		retVal, err := strconv.ParseInt(configValue, 10, 8)
		if err != nil {
			return defaultValue
		}
		return int8(retVal)
	}
	return defaultValue
}

// GetString implements IConfigProvider.
func (e *EnvVarProvider) StringOrDefault(key string, defaultValue string) string {
	if configValue, keyExists := e.config[key]; keyExists {
		return configValue
	}
	return defaultValue
}

// GetUInt implements IConfigProvider.
func (e *EnvVarProvider) UIntOrDefault(key string, defaultValue uint) uint {
	if configValue, keyExists := e.config[key]; keyExists {
		retVal, err := strconv.ParseUint(configValue, 10, bits.UintSize)
		if err != nil {
			return defaultValue
		}
		return uint(retVal)
	}
	return defaultValue
}

// GetUInt16 implements IConfigProvider.
func (e *EnvVarProvider) UInt16OrDefault(key string, defaultValue uint16) uint16 {
	if configValue, keyExists := e.config[key]; keyExists {
		retVal, err := strconv.ParseUint(configValue, 10, 16)
		if err != nil {
			return defaultValue
		}
		return uint16(retVal)
	}
	return defaultValue
}

// GetUInt32 implements IConfigProvider.
func (e *EnvVarProvider) UInt32OrDefault(key string, defaultValue uint32) uint32 {
	if configValue, keyExists := e.config[key]; keyExists {
		retVal, err := strconv.ParseUint(configValue, 10, 32)
		if err != nil {
			return defaultValue
		}
		return uint32(retVal)
	}
	return defaultValue
}

// GetUInt64 implements IConfigProvider.
func (e *EnvVarProvider) UInt64OrDefault(key string, defaultValue uint64) uint64 {
	if configValue, keyExists := e.config[key]; keyExists {
		retVal, err := strconv.ParseUint(configValue, 10, 64)
		if err != nil {
			return defaultValue
		}
		return uint64(retVal)
	}
	return defaultValue
}

// GetUInt8 implements IConfigProvider.
func (e *EnvVarProvider) UInt8OrDefault(key string, defaultValue uint8) uint8 {
	if configValue, keyExists := e.config[key]; keyExists {
		retVal, err := strconv.ParseUint(configValue, 10, 8)
		if err != nil {
			return defaultValue
		}
		return uint8(retVal)
	}
	return defaultValue
}

// SetValue implements IConfigProvider.
func (e *EnvVarProvider) SetValue(key string, value string) {
	e.config[key] = value
}

func NewEnvVarProvider() IConfigProvider {
	evp := &EnvVarProvider{
		config: map[string]string{},
	}
	for _, env := range os.Environ() {
		key, val, _ := strings.Cut(env, "=")
		evp.config[key] = val
	}
	return evp
}
