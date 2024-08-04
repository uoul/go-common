package log

import (
	"fmt"
)

type LogLevel int

const (
	OFF     LogLevel = iota
	FATAL   LogLevel = iota
	ERROR   LogLevel = iota
	WARNING LogLevel = iota
	INFO    LogLevel = iota
	DEBUG   LogLevel = iota
	TRACE   LogLevel = iota
)

func StringToLogLevel(logLevel string) (LogLevel, error) {
	lvl, exists := logLevelMap[logLevel]
	if !exists {
		return INFO, fmt.Errorf("loglevel with name %s does not exist", logLevel)
	}
	return lvl, nil
}

var logLevelMap map[string]LogLevel = map[string]LogLevel{
	"OFF":     OFF,
	"FATAL":   FATAL,
	"ERROR":   ERROR,
	"WARNING": WARNING,
	"INFO":    INFO,
	"DEBUG":   DEBUG,
	"TRACE":   TRACE,
}
