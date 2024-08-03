package log

import (
	"fmt"
	"os"
)

type ConsoleLogger struct {
	level LogLevel
	out   *os.File
	err   *os.File
}

func (l *ConsoleLogger) Trace(message string) {
	if l.level >= TRACE {
		fmt.Fprintf(l.out, "[%s | TRACE] %s\n", currentTimestamp(), message)
	}
}

func (l *ConsoleLogger) Debug(message string) {
	if l.level >= DEBUG {
		fmt.Fprintf(l.out, "[%s | DEBUG] %s\n", currentTimestamp(), message)
	}
}

func (l *ConsoleLogger) Info(message string) {
	if l.level >= INFO {
		fmt.Fprintf(l.out, "[%s | INFO] %s\n", currentTimestamp(), message)
	}
}

func (l *ConsoleLogger) Warning(message string) {
	if l.level >= WARNING {
		fmt.Fprintf(l.out, "[%s | WARNING] %s\n", currentTimestamp(), message)
	}
}

func (l *ConsoleLogger) Error(message string) {
	if l.level >= ERROR {
		fmt.Fprintf(l.err, "[%s | ERROR] %s\n", currentTimestamp(), message)
	}
}

func (l *ConsoleLogger) Fatal(message string) {
	if l.level >= FATAL {
		fmt.Fprintf(l.err, "[%s | FATAL] %s\n", currentTimestamp(), message)
	}
}

func (l *ConsoleLogger) Tracef(format string, a ...any) {
	if l.level >= TRACE {
		fmt.Fprintf(l.out, "[%s | TRACE] %s\n", currentTimestamp(), fmt.Sprintf(format, a...))
	}
}
func (l *ConsoleLogger) Debugf(format string, a ...any) {
	if l.level >= DEBUG {
		fmt.Fprintf(l.out, "[%s | DEBUG] %s\n", currentTimestamp(), fmt.Sprintf(format, a...))
	}
}
func (l *ConsoleLogger) Infof(format string, a ...any) {
	if l.level >= INFO {
		fmt.Fprintf(l.out, "[%s | INFO] %s\n", currentTimestamp(), fmt.Sprintf(format, a...))
	}
}
func (l *ConsoleLogger) Warningf(format string, a ...any) {
	if l.level >= WARNING {
		fmt.Fprintf(l.out, "[%s | WARNING] %s\n", currentTimestamp(), fmt.Sprintf(format, a...))
	}
}
func (l *ConsoleLogger) Errorf(format string, a ...any) {
	if l.level >= ERROR {
		fmt.Fprintf(l.err, "[%s | ERROR] %s\n", currentTimestamp(), fmt.Sprintf(format, a...))
	}
}
func (l *ConsoleLogger) Fatalf(format string, a ...any) {
	if l.level >= FATAL {
		fmt.Fprintf(l.err, "[%s | FATAL] %s\n", currentTimestamp(), fmt.Sprintf(format, a...))
	}
}

func NewConsoleLogger(logLevel LogLevel) ILogger {
	return &ConsoleLogger{
		level: logLevel,
		out:   os.Stdout,
		err:   os.Stderr,
	}
}
