//go:build windows
// +build windows

package log

import (
	"fmt"

	"golang.org/x/sys/windows/svc/debug"
)

type WinEventLogger struct {
	logLevel LogLevel
	elog     debug.Log
}

// Debug implements ILogger.
func (w *WinEventLogger) Debug(message string) {
	if w.logLevel >= DEBUG {
		w.elog.Info(1001, fmt.Sprintf("[%s | DEBUG] %s\n", currentTimestamp(), message))
	}
}

// Debugf implements ILogger.
func (w *WinEventLogger) Debugf(format string, a ...any) {
	if w.logLevel >= DEBUG {
		w.elog.Info(1001, fmt.Sprintf("[%s | DEBUG] %s\n", currentTimestamp(), fmt.Sprintf(format, a...)))
	}
}

// Error implements ILogger.
func (w *WinEventLogger) Error(message string) {
	if w.logLevel >= ERROR {
		w.elog.Error(1001, fmt.Sprintf("[%s | ERROR] %s\n", currentTimestamp(), message))
	}
}

// Errorf implements ILogger.
func (w *WinEventLogger) Errorf(format string, a ...any) {
	if w.logLevel >= ERROR {
		w.elog.Error(1001, fmt.Sprintf("[%s | ERROR] %s\n", currentTimestamp(), fmt.Sprintf(format, a...)))
	}
}

// Fatal implements ILogger.
func (w *WinEventLogger) Fatal(message string) {
	if w.logLevel >= FATAL {
		w.elog.Error(1001, fmt.Sprintf("[%s | FATAL] %s\n", currentTimestamp(), message))
	}
}

// Fatalf implements ILogger.
func (w *WinEventLogger) Fatalf(format string, a ...any) {
	if w.logLevel >= FATAL {
		w.elog.Error(1001, fmt.Sprintf("[%s | FATAL] %s\n", currentTimestamp(), fmt.Sprintf(format, a...)))
	}
}

// Info implements ILogger.
func (w *WinEventLogger) Info(message string) {
	if w.logLevel >= INFO {
		w.elog.Info(1001, fmt.Sprintf("[%s | INFO] %s\n", currentTimestamp(), message))
	}
}

// Infof implements ILogger.
func (w *WinEventLogger) Infof(format string, a ...any) {
	if w.logLevel >= INFO {
		w.elog.Info(1001, fmt.Sprintf("[%s | INFO] %s\n", currentTimestamp(), fmt.Sprintf(format, a...)))
	}
}

// Trace implements ILogger.
func (w *WinEventLogger) Trace(message string) {
	if w.logLevel >= TRACE {
		w.elog.Info(1001, fmt.Sprintf("[%s | TRACE] %s\n", currentTimestamp(), message))
	}
}

// Tracef implements ILogger.
func (w *WinEventLogger) Tracef(format string, a ...any) {
	if w.logLevel >= TRACE {
		w.elog.Info(1001, fmt.Sprintf("[%s | TRACE] %s\n", currentTimestamp(), fmt.Sprintf(format, a...)))
	}
}

// Warning implements ILogger.
func (w *WinEventLogger) Warning(message string) {
	if w.logLevel >= WARNING {
		w.elog.Warning(1001, fmt.Sprintf("[%s | WARNING] %s\n", currentTimestamp(), message))
	}
}

// Warningf implements ILogger.
func (w *WinEventLogger) Warningf(format string, a ...any) {
	if w.logLevel >= WARNING {
		w.elog.Warning(1001, fmt.Sprintf("[%s | WARNING] %s\n", currentTimestamp(), fmt.Sprintf(format, a...)))
	}
}

func NewWinEventLogger(logLevel LogLevel, elog debug.Log) ILogger {
	return &WinEventLogger{
		logLevel: logLevel,
		elog:     elog,
	}
}
