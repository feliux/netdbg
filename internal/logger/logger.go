package logger

import (
	"io"
	"log/slog"
	"sync"
)

// once ensures the logger is initialized only once.
// logger holds the singleton logger instance.
var (
	once   sync.Once
	logger *slog.Logger
)

// InitLogger initializes the logger with a specific configuration.
// Accepts a log level, an io.Writer for output, and a format ("json" or "text").
func InitLogger(level slog.Level, output io.Writer, format string) {
	once.Do(func() {
		var handler slog.Handler

		// Choose the handler based on the format
		switch format {
		case FormatJSON:
			handler = slog.NewJSONHandler(output, &slog.HandlerOptions{Level: level})
		case FormatText:
			handler = slog.NewTextHandler(output, &slog.HandlerOptions{Level: level})
		default:
			panic("unsupported log format: " + format)
		}

		// Create the logger
		logger = slog.New(handler)
	})
}

// GetLogger returns the singleton logger instance.
// If the logger is not initialized, it will panic.
func GetLogger() *slog.Logger {
	if logger == nil {
		panic("logger not initialized. Call InitLogger first.")
	}
	return logger
}

// Info logs an informational message.
func Info(msg string, keysAndValues ...any) {
	GetLogger().Info(msg, keysAndValues...)
}

// Error logs an error message.
func Error(msg string, keysAndValues ...any) {
	GetLogger().Error(msg, keysAndValues...)
}

// Debug logs a debug message.
func Debug(msg string, keysAndValues ...any) {
	GetLogger().Debug(msg, keysAndValues...)
}
