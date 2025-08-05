package servemux

import "github.com/uoul/go-common/serialization"

// -----------------------------------------------------------------------------------------------------------
// Types
// -----------------------------------------------------------------------------------------------------------

type HandlerConfig struct {
	serializer serialization.ISerializer
	maxMemSize int64
}

// -----------------------------------------------------------------------------------------------------------
// Options
// -----------------------------------------------------------------------------------------------------------

func WithHandlerConfigSerializer(serializer serialization.ISerializer) func(*HandlerConfig) {
	return func(hc *HandlerConfig) {
		hc.serializer = serializer
	}
}

func WithHandlerMaxMemSize(size int64) func(*HandlerConfig) {
	return func(hc *HandlerConfig) {
		hc.maxMemSize = size
	}
}

// -----------------------------------------------------------------------------------------------------------
// Constructor
// -----------------------------------------------------------------------------------------------------------

func NewHandlerConfig(opts ...func(*HandlerConfig)) *HandlerConfig {
	c := &HandlerConfig{
		serializer: &serialization.JsonSerializer{},
		maxMemSize: 32 << 20, // 20MB
	}
	for _, o := range opts {
		o(c)
	}
	return c
}
