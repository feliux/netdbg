package logger

import (
	"fmt"
	"log/slog"
	"strings"
)

const (
	// FormatJSON is the JSON log output format.
	FormatJSON = "json"
	// FormatText is the plain text log output format.
	FormatText = "text"
)

// NormalizeFormat validates and normalizes a log format string.
func NormalizeFormat(format string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(format))
	switch normalized {
	case FormatJSON, FormatText:
		return normalized, nil
	default:
		return "", fmt.Errorf("unsupported log format: %s", format)
	}
}

// ParseLevel converts a log level string into a slog.Level.
func ParseLevel(level string) (slog.Level, error) {
	normalized := strings.ToLower(strings.TrimSpace(level))
	switch normalized {
	case "debug":
		return slog.LevelDebug, nil
	case "info", "":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("unsupported log level: %s", level)
	}
}
