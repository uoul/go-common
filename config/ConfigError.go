package config

import "fmt"

type ConfigError struct {
	msg string
}

// Error implements error.
func (c *ConfigError) Error() string {
	return fmt.Sprintf("ConfigError | %s", c.msg)
}

func NewConfigError(format string, a ...any) error {
	return &ConfigError{
		msg: fmt.Sprintf(format, a...),
	}
}
