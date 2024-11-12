package log

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

func StringToLogLevel(logLevel string, defaultLogLevel LogLevel) LogLevel {
	lvl, exists := logLevelMap[logLevel]
	if !exists {
		return defaultLogLevel
	}
	return lvl
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
