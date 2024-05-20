package log

import (
	"flag"
	"io"
	stdlog "log"
	"os"
	"sync"

	"github.com/rs/zerolog"
)

type zlWrapper struct {
	mu      sync.Mutex
	appName string
	logger  *zerolog.Logger
}

var globalLogger Logger
var initOnce sync.Once

// Logger provides common logging functionalities
type Logger interface {
	// BaseLogger returns the underlying logger implementation
	BaseLogger() *zerolog.Logger

	// SetOutput updates the current logger's output to specified io.Writer
	// It also overrides the built-in logger with the updated logger
	SetOutput(w io.Writer)

	// Debug takes formatted message string and logs at Debug level
	Debug(format string, args ...any)

	// Info takes formatted message string and logs at Info level
	Info(format string, args ...any)

	// Warn takes formatted message string and logs at Warn level
	Warn(format string, args ...any)

	// Error takes formatted message string and logs at Error level
	Error(format string, args ...any)

	// Panic takes formatted message string and logs at Panic level. It then calls panic()
	Panic(format string, args ...any)

	// Fatal takes formatted message string and logs at Fatal level. It then calls os.Exit(1)
	Fatal(format string, args ...any)

	// With returns a new wrapped logger with additional context provided by a set
	With(keyValues ...any) Logger
}

// InitGlobalLogger creates a new logger with specified appName and log level
// It also sets the new logger as built in logger
// This should be called prior to calling any of the log methods
func InitGlobalLogger(appName string, level string) {
	initOnce.Do(func() {
		// check if debug flag is already defined. if not, define it
		debugFlag := flag.Lookup("debug")
		if debugFlag == nil {
			flag.Bool("debug", false, "debug mode")
		}
		debug := flag.Lookup("debug").Value.(flag.Getter).Get().(bool)

		level, err := zerolog.ParseLevel(level)
		if err != nil {
			level = zerolog.InfoLevel
		}
		zerolog.SetGlobalLevel(level)
		if debug {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}

		zerolog.LevelFieldName = "l"
		zerolog.MessageFieldName = "m"
		zerolog.TimestampFieldName = "t"
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

		zl := zerolog.New(os.Stderr).With().
			Timestamp().Str("a", appName).Logger()

		globalLogger = &zlWrapper{
			appName: appName,
			logger:  &zl,
		}
		stdlog.SetFlags(0)
		stdlog.SetOutput(&zl)
	})
}

// BaseLogger returns the underlying zerolog logger
func (w *zlWrapper) BaseLogger() *zerolog.Logger {
	return w.logger
}

// SetOutput updates the current logger's output to specified io.Writer
// It also overrides the built-in logger with the updated logger
func (w *zlWrapper) SetOutput(writer io.Writer) {
	w.mu.Lock()
	defer w.mu.Unlock()

	updatedLogger := w.logger.Output(writer)
	w.logger = &updatedLogger
	stdlog.SetFlags(0)
	stdlog.SetOutput(&updatedLogger)
}

// Debug takes formatted message string and logs at Debug level
func (w *zlWrapper) Debug(format string, args ...any) {
	w.logger.Debug().Msgf(format, args...)
}

// Info takes formatted message string and logs at Info level
func (w *zlWrapper) Info(format string, args ...any) {
	w.logger.Info().Msgf(format, args...)
}

// Warn takes formatted message string and logs at Warn level
func (w *zlWrapper) Warn(format string, args ...any) {
	w.logger.Warn().Msgf(format, args...)
}

// Error takes formatted message string and logs at Error level
func (w *zlWrapper) Error(format string, args ...any) {
	w.logger.Error().Msgf(format, args...)
}

// Panic takes formatted message string and logs at Panic level. It then calls panic()
func (w *zlWrapper) Panic(format string, args ...any) {
	w.logger.Panic().Msgf(format, args...)
}

// Fatal takes formatted message string and logs at Fatal level. It then calls os.Exit(1)
func (w *zlWrapper) Fatal(format string, args ...any) {
	w.logger.Fatal().Msgf(format, args...)
}

// With returns a new wrapped logger with additional context
// provided by a set of key/value tuples
func (w *zlWrapper) With(keyValues ...any) Logger {
	zl := w.logger.With().Fields(toFields(keyValues...)).Logger()
	return &zlWrapper{logger: &zl}
}

func toFields(keyValues ...any) map[string]any {
	if len(keyValues)%2 != 0 {
		return make(map[string]any)
	}

	fields := make(map[string]any)
	for i := 0; i < len(keyValues); i += 2 {
		fields[keyValues[i].(string)] = keyValues[i+1]
	}

	return fields
}
