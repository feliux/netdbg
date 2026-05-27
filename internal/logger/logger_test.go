package logger

import (
	"bytes"
	"log/slog"
	"os"
	"testing"
)

// TestLoggerInitialization verifies logger output is written as expected.
func TestLoggerInitialization(t *testing.T) {
	var buf bytes.Buffer

	// Create a new handler with the output redirected to the buffer.
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger = slog.New(handler) // Override the global logger for tests

	logger.Info("Test message", "key", "value")

	if !bytes.Contains(buf.Bytes(), []byte("Test message")) {
		t.Errorf("Expected log message not found")
	}
}

// TestSingletonLogger verifies the logger behaves as a singleton.
func TestSingletonLogger(t *testing.T) {
	InitLogger(slog.LevelInfo, os.Stdout, "text")
	logger1 := GetLogger()
	logger2 := GetLogger()

	if logger1 != logger2 {
		t.Errorf("Logger is not a singleton")
	}
}
