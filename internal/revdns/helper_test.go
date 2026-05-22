package revdns

import (
	"bytes"
	"log/slog"

	"github.com/feliux/netdbg/internal/logger"
)

// setupLogger ensures logger is initialized for each test, improving isolation and reliability.
func setupLogger() {
	var buf bytes.Buffer
	logger.InitLogger(slog.LevelInfo, &buf, "text")
}
