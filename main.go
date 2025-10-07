package main

import (
	"log/slog"
	"os"

	"github.com/feliux/netdbg/cmd"
	"github.com/feliux/netdbg/internal/logger"
)

func main() {
	// Initialize the logger
	logger.InitLogger(slog.LevelInfo, os.Stdout, "json")

	// Execute the CLI
	cmd.Execute()
}
