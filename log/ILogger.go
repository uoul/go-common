package log

type ILogger interface {
	Trace(message string)
	Debug(message string)
	Info(message string)
	Warning(message string)
	Error(message string)
	Fatal(message string)
	Tracef(format string, a ...any)
	Debugf(format string, a ...any)
	Infof(format string, a ...any)
	Warningf(format string, a ...any)
	Errorf(format string, a ...any)
	Fatalf(format string, a ...any)
}
