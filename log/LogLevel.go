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
