package log

import (
	"io"
	"sync"
)

var loggerMutex sync.Mutex

func getLogger() Logger {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	if globalLogger == nil {
		// the logger is not initialized
		// initialize logger with 'un-configured' app name and 'info' log level
		InitGlobalLogger("un-configured", "info")
	}

	return globalLogger
}

// Debug takes formatted message string and logs at Debug level
func Debug(format string, args ...any) {
	getLogger().Debug(format, args...)
}

// Info takes formatted message string and logs at Info level
func Info(format string, args ...any) {
	getLogger().Info(format, args...)
}

// Warn takes formatted message string and logs at Warn level
func Warn(format string, args ...any) {
	getLogger().Warn(format, args...)
}

// Error takes formatted message string and logs at Error level
func Error(format string, args ...any) {
	getLogger().Error(format, args...)
}

// Panic takes formatted message string, logs at Panic level and calls panic()
func Panic(format string, args ...any) {
	getLogger().Panic(format, args...)
}

// Fatal takes formatted message string, logs at Fatal level and then calls os.Exit(1)
func Fatal(format string, args ...any) {
	getLogger().Fatal(format, args...)
}

// SetOutput updates the current logger's output to specified io.Writer
func SetOutput(w io.Writer) {
	getLogger().SetOutput(w)
}

// With returns a new logger with additional context provided by a set of key/value tuples
func With(keyValues ...any) Logger {
	return getLogger().With(keyValues...)
}
